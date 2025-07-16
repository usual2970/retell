import { getTTS, TTSProvider, withModel, withVoice } from "@/pkg/tts";
import {
  Essay,
  EssayCreateReq,
  EssayInfoResp,
  initEssay,
  trans2EssayInfoResp,
} from "@/domain/essay";
import { EssayService } from "../infrastructure/rest/api/essay";

export interface EssayRepository {
  create(essay: Essay): Promise<Essay>;
}

export class Service implements EssayService {
  constructor(private readonly essayRepository: EssayRepository) {}

  async create(essayData: EssayCreateReq): Promise<EssayInfoResp> {
    const newEssay = initEssay(essayData);

    const savedEssay = await this.essayRepository.create(newEssay);

    // 异步转换文本到语音
    this.transTextToSpeech(savedEssay);

    return trans2EssayInfoResp(savedEssay);
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

    // 处理成功的情况
    console.log("TTS audio URL:", resp.audioUrl);
  }
}
