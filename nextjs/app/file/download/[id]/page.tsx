import { Metadata } from "next";

export const metadata: Metadata = {
  title: "File Download",
};

export default async function Page(
  {params}: {params: Promise<{ id: string }>}
) {
  const { id } = await params;
  const fileUrl = `https://shorty.rest/i/t/${id}`;

  return (
    <div className="flex flex-col bg-white rounded-md shadow-lg p-4">
        <b className="mb-2">Thank you for using Shorty!</b>
        <a href={fileUrl} 
          target="_self" 
          className="font-medium text-blue-600 underline dark:text-blue-500 hover:no-underline"
        >
          To download the file, use this link
        </a>
    </div>
  );
}