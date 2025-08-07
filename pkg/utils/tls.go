package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

// IsValidTLSPair 验证 PEM 格式的 cert 和 key 是否有效且匹配
func IsValidTLSPair(domain string, certPEM, keyPEM []byte) (bool, int, error) {
	// 1. 校验输入
	if len(certPEM) == 0 || len(keyPEM) == 0 {
		return false, 0, fmt.Errorf("PEM Certificate or Key is empty")
	}

	// 2. 尝试加载 key pair（这一步同时校验私钥和证书是否匹配）
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return false, 0, fmt.Errorf("x509 key pair load error: %w", err)
	}

	// 3. 校验是否成功提取出证书
	if len(cert.Certificate) == 0 {
		return false, 0, fmt.Errorf("no certificate data found")
	}

	// 4. 解析 DER 格式证
	parsedCert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return false, 0, fmt.Errorf("certificate parse error: %w", err)
	}

	// 校验证书是否与 domain 匹配
	if err := parsedCert.VerifyHostname(domain); err != nil {
		return false, 0, fmt.Errorf("domain mismatch: %w", err)
	}

	// 校验证书是否过期或未生效
	now := time.Now()
	if now.Before(parsedCert.NotBefore) {
		return false, 0, fmt.Errorf("certificate not valid yet (valid from %s)", parsedCert.NotBefore)
	}
	if now.After(parsedCert.NotAfter) {
		return false, 0, fmt.Errorf("certificate has expired (expired at %s)", parsedCert.NotAfter)
	}

	// 计算还有多少天过期
	validDays := int(parsedCert.NotAfter.Sub(now).Hours() / 24)

	return true, validDays, nil
}
