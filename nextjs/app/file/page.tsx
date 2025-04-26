import { Metadata } from "next";
import Form from "next/form";
import { forbidden, unauthorized } from "next/navigation";

export const metadata: Metadata = {
  title: "File",
};

export default function Page() {
  const id = "aaa";
  const captchaBase64 = `data:image/jpeg;base64, ${''}`;
  return (
    <div className="flex flex-col bg-white rounded-md overflow-hidden shadow-xl w-[350px]">
        <div className="flex w-full bg-gradient-to-t from-sky-700 via-sky-800 to-sky-700">
            <p className="text-white ml-2 text-sm">File</p>
        </div>
        <Form action="/file" formMethod="POST" formEncType="multipart/form-data">
            <div className="bg-white rounded-md p-4">
                <input type="hidden" name="id" value={ id }/>
                <input type="file" name="file" accept="*" className="rounded-md mb-2 border-2 border-solid border-gray-400" required/>
                <div className="flex flex-row justify-between items-start">
                    <button id="fileinput" className="flex p-1 pl-2 pr-2 text-center text-white rounded-md shadow-sm bg-sky-800 active:bg-sky-600 hover:bg-sky-700 transition-all">Upload File</button>
                    <div className="flex flex-row rounded-md border border-gray-300">
                        <img className="w-24 h-12 border-r border-gray-300" src={captchaBase64} alt="token"/>
                        <input type="text" inputMode="numeric" placeholder="Captcha..." name="token" className="w-20 h-12 text-center" required/>
                    </div>
                </div>
            </div>
        </Form>
        <div className="mb-1 w-full h-[1px] bg-black"/>
        <p className="p-1">20MB max</p>
    </div>
  );
}
