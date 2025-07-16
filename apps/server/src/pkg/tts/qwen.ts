import { TTS, TTSOption, TTSResponse } from "@/pkg/tts/index";
import { http } from "@/pkg/http";
import { createContextLogger } from "@/pkg/logger";

const url =
  "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation";

type QWENResponse = {
  output: {
    finish_reason: string;
    audio: {
      expires_at: number;
      data: string;
      id: string;
      url: string;
    };
  };
};

export class QWEN implements TTS {
  private logger = createContextLogger("QWEN-TTS");

  async transTextToSpeech(
    text: string,
    ...options: TTSOption[]
  ): Promise<TTSResponse> {
    const requestId = `tts-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    this.logger.info("TTS request started", {
      requestId,
      textLength: text.length,
      hasOptions: options.length > 0,
    });

    const config = {
      model: "qwen-tts",
      voice: "Chelsie",
    };

    options.forEach((option) => option(config));

    this.logger.debug("TTS configuration", {
      requestId,
      config,
    });

    try {
      this.logger.debug("Sending TTS request to QWEN API", {
        requestId,
        url,
        model: config.model,
        voice: config.voice,
      });

      const resp = await http.post<QWENResponse>(
        url,
        {
          model: config.model,
          input: {
            text,
          },
          voice: config.voice,
        },
        {
          headers: {
            Authorization: `Bearer ${process.env.QWEN_API_KEY}`,
          },
        }
      );

      this.logger.info("TTS API response received", {
        requestId,
        status: resp.status,
        finishReason: resp.data.output?.finish_reason,
      });

      if (resp.status !== 200) {
        const error = `HTTP error: ${resp.status}`;
        this.logger.error("TTS request failed with HTTP error", {
          requestId,
          status: resp.status,
          error,
        });

        return {
          error,
          audioUrl: null,
        };
      }

      if (resp.data.output.finish_reason !== "stop") {
        const error = resp.data.output.finish_reason;
        this.logger.warn("TTS request completed with non-stop reason", {
          requestId,
          finishReason: error,
          audioUrl: resp.data.output.audio?.url,
        });

        return {
          error,
          audioUrl: null,
        };
      }

      const audioUrl = resp.data.output.audio.url;
      this.logger.info("TTS request completed successfully", {
        requestId,
        audioUrl,
        audioId: resp.data.output.audio.id,
        expiresAt: new Date(
          resp.data.output.audio.expires_at * 1000
        ).toISOString(),
      });

      return {
        error: null,
        audioUrl,
      };
    } catch (e: any) {
      this.logger.error("TTS request failed with exception", {
        requestId,
        error: {
          name: e.name,
          message: e.message,
          stack: e.stack,
        },
        text: text.substring(0, 100), // 只记录前100个字符
      });

      return {
        error: e.message,
        audioUrl: null,
      };
    }
  }
}
