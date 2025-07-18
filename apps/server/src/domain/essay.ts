import z from "zod";
import { ModelBase } from "./base";

export type Essay = {
  id: number;
  content: string;
  essayType: string;
  file: string;
  fileId: string;
  sentences?: string[];
  taskId: string;
  telegraph: string;
  thumb: string;
  title: string;
  videoLink: string;
} & ModelBase;

export const essayCreateSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().min(1, "Content is required"),
});

export type EssayCreateReq = z.infer<typeof essayCreateSchema>;

export type EssayInfoResp = {
  id: number;
  title: string;
  content: string;
  essayType: string;
  file: string;
  fileId: string;
  sentences?: string[];
  taskId: string;
  telegraph: string;
  thumb: string;
  videoLink: string;
  createdAt?: Date;
  updatedAt?: Date;
};

export const trans2EssayInfoResp = (essay: Essay): EssayInfoResp => {
  return {
    id: essay.id,
    title: essay.title,
    content: essay.content,
    essayType: essay.essayType,
    file: essay.file,
    fileId: essay.fileId,
    sentences: essay.sentences,
    taskId: essay.taskId,
    telegraph: essay.telegraph,
    thumb: essay.thumb,
    videoLink: essay.videoLink,
    createdAt: essay.createdAt,
    updatedAt: essay.updatedAt,
  };
};

export const initEssay = (data: EssayCreateReq): Essay => {
  return {
    id: 0, // Placeholder, will be set by the database
    content: data.content,
    essayType: "default", // Default value, can be changed later
    file: "",
    fileId: "",
    sentences: [],
    taskId: "",
    telegraph: "",
    thumb: "",
    title: data.title,
    videoLink: "",
    createdAt: new Date(),
    updatedAt: new Date(),
  };
};
