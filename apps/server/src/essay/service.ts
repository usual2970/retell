import { getTTS, TTSProvider, withModel, withVoice } from "@/pkg/tts";
import {
  Essay,
  EssayCreateReq,
  EssayInfoResp,
  initEssay,
  trans2EssayInfoResp,
} from "@/domain/essay";
import { EssayService } from "../infrastructure/rest/api/essay";
import { upload } from "@/pkg/upload";
import { search } from "@/pkg/picture";

export interface EssayRepository {
  create(essay: Essay): Promise<Essay>;
  update(id: number, data: Partial<Essay>): Promise<Essay | null>;
}

export class Service implements EssayService {
  constructor(private readonly essayRepository: EssayRepository) {}

  async create(essayData: EssayCreateReq): Promise<EssayInfoResp> {
    const newEssay = initEssay(essayData);

    const savedEssay = await this.essayRepository.create(newEssay);

    // 异步转换文本到语音
    this.transTextToSpeech(savedEssay);

    this.autoSetThumb(savedEssay);

    return trans2EssayInfoResp(savedEssay);
  }

  private async autoSetThumb(essay: Essay): Promise<void> {
    const picture = await search(essay.title);
    if (!picture) return;

    const thumbBuffer = await this.downloadFile(picture.urls.small);

    const uploadResult = await upload(
      thumbBuffer,
      this.generateKeyFromUrl(picture.urls.small, "thumb", picture.format)
    );

    await this.essayRepository.update(essay.id, {
      thumb: uploadResult.url,
    });
  }

  private async transTextToSpeech(essay: Essay): Promise<void> {
    const tts = getTTS(TTSProvider.QWEN);

    const resp = await tts.transTextToSpeech(
      essay.content,
      withModel("qwen-tts"),
      withVoice("Chelsie")
    );

    if (resp.error) {
      // 处理错误
      console.error("TTS error:", resp.error);
      return;
    }

    // TTS 返回的是 audioUrl，下载并上传到 S3，然后保存到数据库
    if (resp.audioUrl) {
      console.log("TTS audio URL:", resp.audioUrl);
      const audioBuffer = await this.downloadFile(resp.audioUrl);
      const uploadResult = await upload(
        audioBuffer,
        this.generateKeyFromUrl(resp.audioUrl)
      );

      console.log("S3 upload result:", uploadResult);

      // 保存到数据库
      await this.essayRepository.update(essay.id, {
        file: uploadResult.url, // 假设 uploadResult 包含上传后的 URL
      });
    }
  }

  private async downloadFile(url: string): Promise<Buffer> {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to download file from ${url}`);
    }
    return Buffer.from(await response.arrayBuffer());
  }

  private generateKeyFromUrl(
    url: string,
    kind: string = "audio",
    format: string = ""
  ): string {
    const urlObj = new URL(url);
    const baseName = urlObj.pathname.split("/").pop() || "";
    urlObj.searchParams;
    // today
    const today =
      new Date().toISOString().split("T")[0]?.replace(/-/g, "") || "";
    return `${kind}/${today}/${baseName}${format ? `.${format}` : ""}`;
  }
}
