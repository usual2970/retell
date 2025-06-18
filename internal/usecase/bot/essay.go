package bot

import (
	"context"
	"ikit-api/internal/domain"
	"ikit-api/internal/util/app"
	"ikit-api/internal/util/audio"
	"ikit-api/internal/util/telegraph"
	"ikit-api/internal/util/zhipu"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

type essayUsecase struct {
	bot *tgbotapi.BotAPI
}

func NewessayUsecase(bot ...*tgbotapi.BotAPI) domain.IessayUsecase {
	if len(bot) > 0 {
		return &essayUsecase{bot: bot[0]}
	}
	return &essayUsecase{}
}

func (e *essayUsecase) UpdateFileId(ctx context.Context, id string, fileId string) error {
	record, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	record.Set("file_id", fileId)
	if err := app.Get().Save(record); err != nil {
		return err
	}

	return nil
}

func (e *essayUsecase) CreateTelegraph(ctx context.Context, id string) error {
	record, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	if record.GetString("content") == "" {
		app.Get().Logger().Info("empty content  no need to upload")
		return nil
	}

	tp := telegraph.New()

	imgUrl := ""
	if record.GetString("thumb") != "" {
		imgUrl = app.Get().Settings().Meta.AppURL + "/api/files/" + record.BaseFilesPath() + "/" + record.GetString("thumb")
	}

	page, err := tp.CreatePage(record.GetString("title"), record.GetString("content"), imgUrl)
	if err != nil {

		app.Get().Logger().Error("create telegraph error:", "err", err)
		return err
	}

	app.Get().Logger().Info("create telegraph success", "page", page)

	// 先重新获取一下record,后面考虑加锁
	cRecord, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	cRecord.Set("telegraph", page.URL)

	if err := app.Get().Save(cRecord); err != nil {
		return err
	}
	app.Get().Logger().Info("create telegraph success", "page", page)
	return nil
}

func (e *essayUsecase) text2Speech(ctx context.Context, id string) error {
	record, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	if record.GetString("content") == "" {
		app.Get().Logger().Info("empty content  no need to upload")
		return nil
	}

	resp, err := audio.Azure(ctx, record.GetString("content"))
	if err != nil {
		return err
	}

	f, err := filesystem.NewFileFromBytes(resp, record.GetString("title"))
	if err != nil {
		return err
	}

	// 先重新获取一下record,后面考虑加锁
	cRecord, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	cRecord.Set("file", []*filesystem.File{f})

	if err := app.Get().Save(cRecord); err != nil {
		return err
	}

	return nil
}

func (e *essayUsecase) Notify(ctx context.Context, req *domain.TtsAsyncResp) error {
	app.Get().Logger().Info("essay notify:", "req", req)

	record, err := app.Get().FindFirstRecordByData("essay", "task_id", req.Data.TaskId)
	if err != nil {
		return err
	}

	record.Set("sentences", req.Data.Sentences)

	f1, err := filesystem.NewFileFromURL(ctx, req.Data.AudioAddress)
	if err != nil {
		return err
	}
	record.Set("file", []*filesystem.File{f1})

	if err := app.Get().Save(record); err != nil {
		return err
	}
	return nil
}

func (e *essayUsecase) Add(ctx context.Context, req *domain.AddessayReq) error {

	collection, err := app.Get().FindCollectionByNameOrId("essay")
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)

	record.Set("title", req.Title)
	record.Set("content", req.Content)

	apiKey := os.Getenv("ZHIPU_API_KEY")
	zp := zhipu.NewZhipu(apiKey)

	url, err := zp.GenerateImg(ctx, req.Title)
	if err != nil {
		return err
	}
	f, _ := filesystem.NewFileFromURL(ctx, url)

	record.Set("thumb", []*filesystem.File{f})

	// 保存到数据库
	if err := app.Get().Save(record); err != nil {
		app.Get().Logger().Error("save essay error:", "err", err)
		return err
	}

	go func() {
		if err := e.CreateTelegraph(context.Background(), record.Id); err != nil {
			app.Get().Logger().Error("createTelegraph error:", "err", err)
		} else {
			app.Get().Logger().Info("success createTelegraph", "id", record.Id)
		}
	}()

	// 文字转换成语音
	go func() {
		if err := e.text2Speech(context.Background(), record.Id); err != nil {
			app.Get().Logger().Error("text2speech error:", "err", err)
		} else {
			app.Get().Logger().Info("success text2speech", "id", record.Id)
		}
	}()

	return nil
}

func (e *essayUsecase) List(ctx context.Context, req *domain.ListessayReq) ([]domain.Essay, error) {

	records, err := app.Get().FindRecordsByFilter("essay", req.Filter, "-created", req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	rs := make([]domain.Essay, 0, len(records))
	for _, record := range records {
		rs = append(rs, domain.Essay{
			Meta: domain.Meta{
				Id:      record.Id,
				Created: record.GetDateTime("created").Time(),
				Updated: record.GetDateTime("updated").Time(),
			},
			Title: record.GetString("title"),
		})
	}

	return rs, nil
}

func (e *essayUsecase) Delete(ctx context.Context, id string) error {
	record, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return err
	}

	return app.Get().Delete(record)
}

func (e *essayUsecase) Detail(ctx context.Context, id string) (*domain.Essay, error) {
	record, err := app.Get().FindRecordById("essay", id)
	if err != nil {
		return nil, err
	}
	file := ""
	if record.GetString("file") != "" {
		file = app.Get().Settings().Meta.AppURL + "/api/files/" + record.BaseFilesPath() + "/" + record.GetString("file")
	}

	thumb := ""
	if record.GetString("thumb") != "" {
		thumb = app.Get().Settings().Meta.AppURL + "/api/files/" + record.BaseFilesPath() + "/" + record.GetString("thumb")
	}
	rs := &domain.Essay{
		Meta: domain.Meta{
			Id:      record.Id,
			Created: record.GetDateTime("created").Time(),
			Updated: record.GetDateTime("updated").Time(),
		},
		Title:     record.GetString("title"),
		Content:   record.GetString("content"),
		File:      file,
		FileId:    record.GetString("file_id"),
		Thumb:     thumb,
		Telegraph: record.GetString("telegraph"),
	}
	return rs, nil
}
