import { Metadata } from "next";
import Form from "next/form";

export const metadata: Metadata = {
  title: "File",
};

export default async function Page(
  {params}: {params: Promise<{ id: string }>}
) {
  const { id } = await params;
  const fileName = "test";
  const fileSize = `0.12 MB`;
  const fileUrl = `https://shorty.rest/i/t/${id}`;
  const fileViewUrl = `https://shorty.rest/i/t/${id}`;
  const captchaid = "aaa";
  const captchaBase64 = `data:image/jpeg;base64, ${''}`;

  return (
    <div className="bg-white rounded-md shadow-lg p-4 min-w-[400px]">
        <div className="flex flex-row justify-between items-start w-full">
            <div className="flex flex-row justify-start items-start">
                <img src="/static/file.png" className="w-12 h-12 mr-2 p-0 object-contain" alt="file"/>
                <div className="flex flex-col justify-between">
                    <p className="mb-1 text-md">{fileName}</p>
                    <p className="mb-2 text-sm">{fileSize}</p>
                </div>
            </div>
            <button className="ml-2 pl-1 pr-1 rounded-md text-white font-bold bg-red-600 hover:bg-red-400 active:bg-red-500">Report</button>
        </div>
        <div className="flex flex-col mb-2 w-full">        
            <Form action={fileUrl} formMethod="GET">
              <div className="bg-white rounded-md mt-2">
                  <input type="hidden" name="id" value={captchaid}/>
                  <div className="flex flex-row justify-between items-start">
                      <button id="notifyanchor" className="flex p-1 pl-2 pr-2 text-center text-white rounded-md shadow-sm bg-sky-800 active:bg-sky-600 hover:bg-sky-700 transition-all">Download File</button>
                      <div className="flex flex-row rounded-md border border-gray-300">
                          <img className="w-24 h-12 border-r border-gray-300" src={captchaBase64} alt="token"/>
                          <input type="text" inputMode="numeric" placeholder="Captcha..." name="token" className="w-20 h-12 text-center" required/>
                      </div>
                      </div>
                  </div>
              </Form>
            </div>
        <p>URL:</p>
        <textarea className="w-full rounded-sm p-1 bg-gray-200 resize-none">{fileViewUrl}</textarea>
    </div>
  );
}