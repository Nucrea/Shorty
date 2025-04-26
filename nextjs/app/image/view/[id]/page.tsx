import api from "@/api/api";
import { Metadata } from "next";
import { notFound } from "next/navigation";

export const metadata: Metadata = {
  title: "Image",
};

interface ImageViewParams {
    ViewUrl: string
    ImageUrl: string
    ThumbnailUrl: string
    FileName: string
    SizeMB: number
}

// func bbCode(p ImageViewParams) string {
//     return fmt.Sprintf("[URL=%s][IMG]%s[/IMG][/URL]", templ.URL(p.ViewUrl), templ.URL(p.ThumbnailUrl))
// }

export default async function Page(
  {params}: {params: Promise<{ id: string }>}
) {
  const { id } = await params;

  const info = await api.GetImageInfo(id);
  if (!info) {
    return notFound();
  }

  const { name, size, originalUrl, thumbnailUrl } = info!
  const imageSize = `${size.toFixed(1)} MB`;
  const bbCode = `[URL=${originalUrl}][IMG]${thumbnailUrl}[/IMG][/URL]`;
  const viewUrl = `http://localhost:3001/image/view/${info.id}`;

  return (
    <div className="bg-white rounded-md shadow-lg p-4">
        <div className="flex flex-row justify-between items-start w-full">
            <div className="flex flex-col">
                <p className="mb-1 text-md">{name}</p>
                <p className="mb-2 text-sm">{imageSize}</p>
            </div>
            <button className="p-1 rounded-md text-white font-bold bg-red-600 hover:bg-red-400 active:bg-red-500">Report</button>
        </div>
        <div className="mb-1 w-full h-[1px] bg-black"/>
        <a href={originalUrl} target="_self">
            <img className="rounded-md w-auto h-[50vh]" src={originalUrl} alt={name}/>
        </a>
        <div className="mb-1 w-full h-[1px] bg-black"/>
        <p>URL:</p>
        <textarea className="w-full rounded-sm p-1 bg-gray-200 resize-none" value={viewUrl} contentEditable={false}/>
        <p className="mt-1">BB-Code:</p>
        <textarea className="w-full rounded-sm p-1 bg-gray-200 resize-none" value={bbCode} contentEditable={false}/>
    </div>
  );
}