import { createApi, OrderBy } from "unsplash-js";

const api = createApi({
  accessKey: process.env.UNSPLASH_ACCESS_KEY || "",
});

type SearchResponse = {
  urls: {
    full: string;
    raw: string;
    regular: string;
    small: string;
    thumb: string;
  };
  format: string;
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

  const fullUrl = photo.urls.full;

  const urlObj = new URL(fullUrl);
  const format = urlObj.searchParams.get("fm");

  return {
    urls: photo.urls,
    format: format || "",
  };
};
