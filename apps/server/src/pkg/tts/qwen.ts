import { TTS, TTSOption } from "@/pkg/tts/index";
import { http } from "@/pkg/http";

const url =
  "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation";

export class QWEN implements TTS {
  async transTextToSpeech(
    text: string,
    ...options: TTSOption[]
  ): Promise<void> {
    const config = {
      model: "qwen-tts",
      voice: "Chelsie",
    };

    options.forEach((option) => option(config));

    http.post(url, {
      model: config.model,
      input: {
        text,
      },
      voice: config.voice,
    });
  }
}
