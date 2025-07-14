import {
  Essay,
  EssayCreateReq,
  EssayInfoResp,
  initEssay,
  trans2EssayInfoResp,
} from "../domain/essay";
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
    

    return trans2EssayInfoResp(savedEssay);
  }
}
