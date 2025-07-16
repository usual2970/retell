import { createApi, OrderBy } from "unsplash-js";

const api = createApi({
  accessKey: process.env.UNSPLASH_ACCESS_KEY || "",
});

type SearchResponse = {
  full: string;
  raw: string;
  regular: string;
  small: string;
  thumb: string;
};

export const search = async (
  query: string
): Promise<SearchResponse | undefined> => {
  const resp = await api.search.getPhotos({
    query: query,
    page: 1,
    perPage: 1,
    orderBy: "relevant",
  });

  const photo = resp?.response?.results[0];
  if (!photo) return;

  return photo.urls as SearchResponse;
};
