package file

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/panjf2000/ants/v2"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/encode"
	"github.com/usual2970/retell/pkg/logger"

	"github.com/aliyun/credentials-go/credentials"
)

type FileRepository interface {
	Save(ctx context.Context, file *domain.File) error
}

type Service struct {
	fileRepo FileRepository
}

func NewService(fileRepo FileRepository) *Service {
	return &Service{
		fileRepo: fileRepo,
	}
}

func (s *Service) PolicyToken(ctx context.Context) (*domain.FilePolicyTokenResp, error) {
	product := "oss"
	conf := config.GetConfig().OSSDirect

	// 设置bucket所处地域
	region := conf.Region
	// 替换为您的bucket名称
	bucketName := conf.Bucket
	// 设置 OSS 上传地址
	host := fmt.Sprintf("https://%s", conf.Domain)
	// 设置上传目录

	now := time.Now()
	dir := fmt.Sprintf("uploads/%d/%02d/%02d/",
		now.Year(),
		now.Month(),
		now.Day(),
	)

	config := new(credentials.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(conf.AccessKeyId).
		SetAccessKeySecret(conf.AccessKeySecret).
		SetRoleArn(conf.RoleArn).
		SetRoleSessionName(conf.RoleSessionName).
		SetPolicy("").
		SetRoleSessionExpiration(3600)

	// 根据配置创建凭证提供器
	provider, err := credentials.NewCredential(config)
	if err != nil {
		log.Fatalf("NewCredential fail, err:%v", err)
	}

	// 从凭证提供器获取凭证
	cred, err := provider.GetCredential()
	if err != nil {
		return nil, constant.NewXError(100, "获取OSS凭证失败", err)
	}

	// 构建policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": bucketName},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
		},
	}

	// 将policy转换为 JSON 格式
	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalf("json.Marshal fail, err:%v", err)
	}

	// 构造待签名字符串（StringToSign）
	stringToSign := base64.StdEncoding.EncodeToString([]byte(policy))

	hmacHash := func() hash.Hash { return sha256.New() }
	// 构建signing key
	signingKey := "aliyun_v4" + *cred.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	// 构建返回给前端的表单
	return &domain.FilePolicyTokenResp{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product),
		Date:             utcTime.UTC().Format("20060102T150405Z"),
		Signature:        signature,
		Host:             host, // 返回 OSS 上传地址
		Dir:              dir,  // 返回上传目录
	}, nil

}

func (s *Service) ValidateSQL(ctx context.Context, req *domain.FileValidateSQLReq) (*domain.FileValidateSQLResp, error) {
	// 0. 检查文件大小限制
	maxFileSize := int64(50 * 1024 * 1024) // 50MB SQL文件限制
	fileSize, err := s.getRemoteFileSize(req.URL)
	if err != nil {
		return nil, fmt.Errorf("获取文件大小失败: %v", err)
	}

	if fileSize > maxFileSize {
		return &domain.FileValidateSQLResp{
			Error: fmt.Sprintf("SQL文件过大，超过限制 %d MB", maxFileSize/(1024*1024)),
		}, nil
	}

	// 1. 流式下载并验证SQL文件
	_, message, err := s.validateSQLStream(ctx, req.URL)

	if err != nil {
		message = err.Error()
	}

	return &domain.FileValidateSQLResp{
		Error: message,
	}, nil
}

// validateSQLStream 流式验证SQL文件
func (s *Service) validateSQLStream(ctx context.Context, url string) (bool, string, error) {
	// 发起HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Minute, // 5分钟超时
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 使用缓冲读取器进行流式处理
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // 64KB初始缓冲区，最大1MB

	var (
		lineNumber       = 0
		inMultiComment   = false
		currentStatement strings.Builder
		dangerousOps     []string
	)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if len(line) == 0 {
			continue
		}

		// 处理多行注释
		if inMultiComment {
			if strings.Contains(line, "*/") {
				inMultiComment = false
				// 处理注释结束后的剩余内容
				parts := strings.Split(line, "*/")
				if len(parts) > 1 {
					line = strings.TrimSpace(parts[1])
					if len(line) == 0 {
						continue
					}
				} else {
					continue
				}
			} else {
				continue
			}
		}

		// 检查多行注释开始
		if strings.Contains(line, "/*") {
			parts := strings.Split(line, "/*")
			line = strings.TrimSpace(parts[0])
			if strings.Contains(parts[1], "*/") {
				// 单行内的多行注释
				commentParts := strings.Split(parts[1], "*/")
				if len(commentParts) > 1 {
					line += " " + strings.TrimSpace(commentParts[1])
				}
			} else {
				inMultiComment = true
			}
		}

		// 跳过单行注释
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			continue
		}

		// 移除行尾注释
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		if len(line) == 0 {
			continue
		}

		// 累积SQL语句（处理跨行语句）
		currentStatement.WriteString(" ")
		currentStatement.WriteString(line)

		// 检查是否是完整的语句（以分号结尾）
		if strings.HasSuffix(line, ";") {
			statement := strings.TrimSpace(currentStatement.String())

			// 检查危险操作
			if dangerous := s.checkDangerousOperation(statement); dangerous != "" {
				dangerousOps = append(dangerousOps, fmt.Sprintf("第%d行: %s", lineNumber, dangerous))
			}

			// 重置语句缓冲区
			currentStatement.Reset()
		}

		// 限制检查的行数，防止无限循环
		if lineNumber > 100000 {
			return false, "SQL文件过大，超过10万行限制", nil
		}
	}

	// 检查最后一个可能没有分号结尾的语句
	if currentStatement.Len() > 0 {
		statement := strings.TrimSpace(currentStatement.String())
		if dangerous := s.checkDangerousOperation(statement); dangerous != "" {
			dangerousOps = append(dangerousOps, fmt.Sprintf("第%d行: %s", lineNumber, dangerous))
		}
	}

	if err := scanner.Err(); err != nil {
		return false, "", fmt.Errorf("读取文件流失败: %v", err)
	}

	// 返回验证结果
	if len(dangerousOps) > 0 {
		message := "检测到危险操作:\n" + strings.Join(dangerousOps, "\n")
		return false, message, nil
	}

	return true, "", nil
}

// checkDangerousOperation 检查SQL语句中的危险操作
func (s *Service) checkDangerousOperation(statement string) string {
	// 转换为大写便于匹配
	upperStatement := strings.ToUpper(statement)

	// 移除多余的空格
	upperStatement = strings.Join(strings.Fields(upperStatement), " ")

	// 定义危险操作模式
	dangerousPatterns := []struct {
		pattern     string
		description string
		regex       *regexp.Regexp
		notRegex    *regexp.Regexp
	}{
		{
			pattern:     "DROP COLUMN",
			description: "删除列操作",
			regex:       regexp.MustCompile(`\bALTER\s+TABLE\s+\w+\s+DROP\s+COLUMN\s+\w+`),
			notRegex:    nil,
		},
		{
			pattern:     "DROP TABLE (without IF EXISTS)",
			description: "直接删除表操作",
			// 匹配 DROP TABLE 但不包含 IF EXISTS 的情况
			regex:    regexp.MustCompile(`\bDROP\s+TABLE\s+.?\w+`),
			notRegex: regexp.MustCompile(`\bDROP\s+TABLE\s+IF\s+EXISTS\s+.?\w+`),
		},
		{
			pattern:     "DROP DATABASE (without IF EXISTS)",
			description: "直接删除数据库操作",
			// 匹配 DROP DATABASE/SCHEMA 但不包含 IF EXISTS 的情况
			regex:    regexp.MustCompile(`\bDROP\s+(DATABASE|SCHEMA)\s+\.?w+`),
			notRegex: regexp.MustCompile(`\bDROP\s+(DATABASE|SCHEMA)\s+IF\s+EXISTS\s+.?\w+`),
		},
	}

	// 检查每个危险模式
	for _, pattern := range dangerousPatterns {
		first := false
		second := false

		if pattern.regex.MatchString(upperStatement) {
			first = true

		}

		if pattern.notRegex != nil && pattern.notRegex.MatchString(upperStatement) {
			second = true
		}

		if first && !second {
			return fmt.Sprintf("%s: %s", pattern.description, s.extractRelevantPart(statement, pattern.pattern))
		}
	}

	// 额外检查：简单的字符串匹配作为备用（保留删除列的检查）
	if strings.Contains(upperStatement, "DROP COLUMN") {
		return "删除列操作: " + s.extractRelevantPart(statement, "DROP COLUMN")
	}

	// 检查直接的 DROP TABLE（不包含 IF EXISTS）
	if strings.Contains(upperStatement, "DROP TABLE") && !strings.Contains(upperStatement, "IF EXISTS") {
		return "直接删除表操作: " + s.extractRelevantPart(statement, "DROP TABLE")
	}

	// 检查直接的 DROP DATABASE（不包含 IF EXISTS）
	if (strings.Contains(upperStatement, "DROP DATABASE") || strings.Contains(upperStatement, "DROP SCHEMA")) &&
		!strings.Contains(upperStatement, "IF EXISTS") {
		return "直接删除数据库操作: " + s.extractRelevantPart(statement, "DROP DATABASE")
	}

	return ""
}

// extractRelevantPart 提取相关的SQL部分用于错误报告
func (s *Service) extractRelevantPart(statement, pattern string) string {
	// 限制返回的SQL片段长度
	maxLength := 100

	upperStatement := strings.ToUpper(statement)
	upperPattern := strings.ToUpper(pattern)

	index := strings.Index(upperStatement, upperPattern)
	if index == -1 {
		// 如果没找到确切匹配，返回整个语句的前100个字符
		if len(statement) > maxLength {
			return statement[:maxLength] + "..."
		}
		return statement
	}

	// 提取包含危险操作的部分
	start := index
	end := index + len(pattern) + 50 // 包含模式后的50个字符

	if start > 20 {
		start -= 20 // 包含模式前的20个字符
	} else {
		start = 0
	}

	if end > len(statement) {
		end = len(statement)
	}

	result := statement[start:end]

	if start > 0 {
		result = "..." + result
	}
	if end < len(statement) {
		result = result + "..."
	}

	return result
}

func (s *Service) Structure(ctx context.Context, url string) (*domain.FileStructureResp, error) {
	// 0. 检查文件大小限制
	maxFileSize := int64(100 * 1024 * 1024) // 100MB
	fileSize, err := s.getRemoteFileSize(url)
	if err != nil {
		return nil, fmt.Errorf("获取文件大小失败: %v", err)
	}

	if fileSize > maxFileSize {
		return nil, fmt.Errorf("文件过大，超过限制 %d MB", maxFileSize/(1024*1024))
	}

	// 1. 下载文件但使用流式处理
	tempFilePath, err := s.downloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %v", err)
	}
	// 确保临时文件被删除
	defer os.Remove(tempFilePath)

	// 2. 解析文件结构（优化版本）
	structure, err := s.parseZipStructureOptimized(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("解析文件结构失败: %v", err)
	}

	return structure, nil
}

// getRemoteFileSize 获取远程文件大小而无需下载整个文件
func (s *Service) getRemoteFileSize(url string) (int64, error) {
	// 发送 HEAD 请求，只获取头信息
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("HTTP HEAD请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	// 从响应头获取文件大小
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return 0, fmt.Errorf("无法获取文件大小")
	}

	return contentLength, nil
}

// downloadFile 下载远程文件到临时目录
func (s *Service) downloadFile(url string) (string, error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "download-*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer tempFile.Close()
	tempFilePath := tempFile.Name()

	// 发起HTTP请求下载文件
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 将响应内容写入临时文件
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("写入临时文件失败: %v", err)
	}

	return tempFilePath, nil
}

// parseZipStructureOptimized 解析ZIP格式文件的结构（内存优化版本）
func (s *Service) parseZipStructureOptimized(filePath string) (*domain.FileStructureResp, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %v", err)
	}
	defer r.Close()

	// 构建文件结构树
	structure := &domain.FileStructureResp{
		Type:  "archive",
		Name:  filepath.Base(filePath),
		Files: []*domain.FileNode{},
		Size:  0,
	}

	// 获取文件总大小
	if fileInfo, err := os.Stat(filePath); err == nil {
		structure.Size = fileInfo.Size()
	}

	// 使用映射表跟踪已创建的目录节点
	// 仅存储路径到节点的引用，不存储完整内容
	dirMap := make(map[string]*domain.FileNode)

	// 限制处理的最大条目数
	maxEntries := 10000
	entryCount := 0

	// 遍历 ZIP 文件中的所有条目
	for _, f := range r.File {
		// 限制最大条目数
		entryCount++
		if entryCount > maxEntries {
			return nil, fmt.Errorf("ZIP文件条目数超出限制 (%d)", maxEntries)
		}

		// 跳过 macOS 特殊目录和隐藏文件
		if strings.HasPrefix(f.Name, "__MACOSX/") || strings.Contains(f.Name, "/.") {
			continue
		}

		// 标准化路径
		path := strings.TrimSuffix(strings.ReplaceAll(f.Name, "\\", "/"), "/")
		if path == "" {
			continue
		}

		// 文件层级过深可能导致内存问题
		if len(strings.Split(path, "/")) > 20 {
			continue // 跳过过深的路径
		}

		// 创建节点（只存储必要信息）
		isDir := f.FileInfo().IsDir()

		// 处理目录
		if isDir {
			// 确保该目录节点已创建
			s.ensureDirectoryExists(path, structure, dirMap)
			continue
		}

		// 处理文件
		dirname := filepath.Dir(path)
		if dirname == "." {
			// 根目录下的文件
			node := &domain.FileNode{
				Name: filepath.Base(path),
				Path: path,
				Type: "file",
				Size: int64(f.UncompressedSize64),
			}
			structure.Files = append(structure.Files, node)
		} else {
			// 确保父目录节点存在
			parent := s.ensureDirectoryExists(dirname, structure, dirMap)

			// 将文件添加到父目录
			node := &domain.FileNode{
				Name: filepath.Base(path),
				Path: path,
				Type: "file",
				Size: int64(f.UncompressedSize64),
			}
			parent.Children = append(parent.Children, node)
		}
	}

	// 计算目录大小
	s.calculateDirectorySize(structure.Files)

	return structure, nil
}

// ensureDirectoryExists 确保指定路径的目录节点存在，如不存在则创建
func (s *Service) ensureDirectoryExists(path string, structure *domain.FileStructureResp, dirMap map[string]*domain.FileNode) *domain.FileNode {
	// 检查目录节点是否已存在
	dirNode, exists := dirMap[path]
	if exists {
		return dirNode
	}

	// 如果是根目录路径
	if path == "." {
		return nil
	}

	// 创建新的目录节点
	dirNode = &domain.FileNode{
		Name:     filepath.Base(path),
		Path:     path,
		Type:     "directory",
		Children: []*domain.FileNode{},
		Size:     0,
	}

	// 添加到目录映射
	dirMap[path] = dirNode

	// 处理父目录
	parent := filepath.Dir(path)
	if parent == "." {
		// 根级目录
		structure.Files = append(structure.Files, dirNode)
	} else {
		// 确保父目录存在并添加到父目录
		parentNode := s.ensureDirectoryExists(parent, structure, dirMap)
		parentNode.Children = append(parentNode.Children, dirNode)
	}

	return dirNode
}

// calculateDirectorySize 递归计算目录大小
func (s *Service) calculateDirectorySize(nodes []*domain.FileNode) int64 {
	var totalSize int64
	for i := range nodes {
		node := nodes[i]
		if node.Type == "directory" {
			node.Size = s.calculateDirectorySize(node.Children)
		}
		totalSize += node.Size
	}
	return totalSize
}

func (s *Service) Upload(ctx context.Context, file *multipart.FileHeader) (*domain.FileUploadResp, error) {
	ossConf := config.GetConfig().OSS

	// 创建 OSSClient 实例
	client, err := oss.New(ossConf.Endpoint, ossConf.AccessKeyId, ossConf.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %v", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(ossConf.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取 Bucket 失败: %v", err)
	}

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	defer src.Close()

	// 生成带日期的文件路径
	now := time.Now()
	uri := encode.Uri()
	objectKey := fmt.Sprintf("uploads/%d/%02d/%02d/%s%s",
		now.Year(),
		now.Month(),
		now.Day(),
		uri,
		filepath.Ext(file.Filename),
	)

	// 检查文件大小，决定使用哪种上传方式
	if file.Size > 10*1024*1024 { // 如果大于10MB，使用分片上传
		// 初始化分片上传
		imur, err := bucket.InitiateMultipartUpload(objectKey)
		if err != nil {
			return nil, fmt.Errorf("初始化分片上传失败: %v", err)
		}

		// 分片大小，5MB
		chunkSize := int64(5 * 1024 * 1024)

		// 计算分片数
		chunks := int((file.Size + chunkSize - 1) / chunkSize)

		// 保存每个分片的ETag
		parts := make([]oss.UploadPart, chunks)

		// 逐个上传分片
		for i := 0; i < chunks; i++ {
			// 计算当前分片大小
			partSize := chunkSize
			if i == chunks-1 {
				partSize = file.Size - int64(i)*chunkSize
			}

			// 创建缓冲区读取当前分片
			buffer := make([]byte, partSize)
			n, err := io.ReadFull(src, buffer)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				// 取消分片上传
				bucket.AbortMultipartUpload(imur)
				return nil, fmt.Errorf("读取文件分片失败: %v", err)
			}

			// 上传当前分片
			part, err := bucket.UploadPart(imur, bytes.NewReader(buffer[:n]), partSize, i+1)
			if err != nil {
				// 取消分片上传
				bucket.AbortMultipartUpload(imur)
				return nil, fmt.Errorf("上传分片失败: %v", err)
			}
			parts[i] = part
		}

		// 完成分片上传
		_, err = bucket.CompleteMultipartUpload(imur, parts)
		if err != nil {
			return nil, fmt.Errorf("完成分片上传失败: %v", err)
		}
	} else {
		// 常规上传但使用缓冲流
		bufferSize := 4 * 1024 * 1024 // 4MB buffer
		buffer := make([]byte, bufferSize)

		// 创建管道实现流式上传
		reader, writer := io.Pipe()

		// 在后台将文件写入管道
		go func() {
			defer writer.Close()

			for {
				n, err := src.Read(buffer)
				if err != nil && err != io.EOF {
					logger.WithFields(map[string]interface{}{
						"error": err,
					}).Error("读取上传文件失败")
					return
				}

				if n == 0 {
					break
				}

				if _, err := writer.Write(buffer[:n]); err != nil {
					logger.WithFields(map[string]interface{}{
						"error": err,
					}).Error("写入管道失败")
					return
				}
			}
		}()

		// 使用管道进行流式上传
		err = bucket.PutObject(objectKey, reader)
		if err != nil {
			return nil, fmt.Errorf("上传到 OSS 失败: %v", err)
		}
	}

	// 生成文件访问 URL
	fileURL := fmt.Sprintf("https://%s/%s", ossConf.Domain, objectKey)

	// 异步保存文件信息到数据库
	ants.Submit(func() {
		file := &domain.File{
			Name: file.Filename,
			Url:  fmt.Sprintf("/%s", objectKey),
			Uri:  uri,
		}
		ctx := context.Background()
		if err := s.fileRepo.Save(ctx, file); err != nil {
			logger.WithFields(map[string]interface{}{
				"fileURL": fileURL,
				"error":   err.Error(),
			}).Warn("保存文件记录失败")
		}
	})

	return &domain.FileUploadResp{
		URL: fileURL,
	}, nil
}

func (s *Service) UploadLocalFile(ctx context.Context, filePath string) (*domain.FileUploadResp, error) {
	ossConf := config.GetConfig().OSS

	// 创建 OSSClient 实例
	client, err := oss.New(ossConf.Endpoint, ossConf.AccessKeyId, ossConf.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %v", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(ossConf.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取 Bucket 失败: %v", err)
	}

	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开本地文件失败: %v", err)
	}
	//defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 生成带日期的文件路径
	now := time.Now()
	uri := encode.Uri()
	objectKey := fmt.Sprintf("uploads/%d/%02d/%02d/%s%s",
		now.Year(),
		now.Month(),
		now.Day(),
		uri,
		filepath.Ext(info.Name()),
	)

	// 判断是否需要使用分片上传
	if info.Size() > 10*1024*1024 {
		// 初始化分片上传
		imur, err := bucket.InitiateMultipartUpload(objectKey)
		if err != nil {
			return nil, fmt.Errorf("初始化分片上传失败: %v", err)
		}

		chunkSize := int64(5 * 1024 * 1024)
		chunks := int((info.Size() + chunkSize - 1) / chunkSize)
		parts := make([]oss.UploadPart, chunks)

		for i := 0; i < chunks; i++ {
			partSize := chunkSize
			if i == chunks-1 {
				partSize = info.Size() - int64(i)*chunkSize
			}

			buffer := make([]byte, partSize)
			n, err := io.ReadFull(file, buffer)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				bucket.AbortMultipartUpload(imur)
				return nil, fmt.Errorf("读取文件分片失败: %v", err)
			}

			part, err := bucket.UploadPart(imur, bytes.NewReader(buffer[:n]), partSize, i+1)
			if err != nil {
				bucket.AbortMultipartUpload(imur)
				return nil, fmt.Errorf("上传分片失败: %v", err)
			}
			parts[i] = part
		}

		_, err = bucket.CompleteMultipartUpload(imur, parts)
		if err != nil {
			return nil, fmt.Errorf("完成分片上传失败: %v", err)
		}
	} else {
		// 小文件直接上传
		err = bucket.PutObject(objectKey, file)
		if err != nil {
			return nil, fmt.Errorf("上传到 OSS 失败: %v", err)
		}
	}

	// 构造 URL
	fileURL := fmt.Sprintf("https://%s/%s", ossConf.Domain, objectKey)

	// 异步保存文件记录
	ants.Submit(func() {
		fileRecord := &domain.File{
			Name: info.Name(),
			Url:  fmt.Sprintf("/%s", objectKey),
			Uri:  uri,
		}
		if err := s.fileRepo.Save(context.Background(), fileRecord); err != nil {
			logger.WithFields(map[string]interface{}{
				"fileURL": fileURL,
				"error":   err.Error(),
			}).Warn("保存文件记录失败")
		}
	})

	file.Close()
	err = os.Remove(filePath)
	if err != nil {
		return nil, err
	}

	return &domain.FileUploadResp{
		URL: fileURL,
	}, nil
}
