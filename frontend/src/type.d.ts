export interface IValues {
  pathFolder: string,
  playlistId: string,
}
export type VideoState = {
  ID: string,
  Author: string,
  Title: string,
  Duration: string,
  Thumbnails: {
    URL: string,
    Width: number,
    Height: number,
  }[],
  Status: "pending" | "downloading" | "done" | "error",
}