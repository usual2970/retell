import { QWEN } from "./qwen";

export enum TTSProvider {
  QWEN = "QWEN",
}

export type TTSConfig = {
  model: string;
  voice: string;
};

export type TTSOption = (config: TTSConfig) => void;

export const withModel =
  (model: string): TTSOption =>
  (config) => {
    config.model = model;
  };

export const withVoice =
  (voice: string): TTSOption =>
  (config) => {
    config.voice = voice;
  };

export interface TTS {
  transTextToSpeech(text: string, ...options: TTSOption[]): Promise<void>;
}

export const getTTS = (provider: TTSProvider): TTS => {
  switch (provider) {
    case TTSProvider.QWEN:
      return new QWEN();
    default:
      throw new Error(`Unsupported TTS provider: ${provider}`);
  }
};
