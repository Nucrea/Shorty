import { Metadata } from "next";
import Form from "next/form";

export const metadata: Metadata = {
  title: "Link",
};

export default function Page() {
  return (
    <div className="flex flex-col justify-between bg-white rounded-lg shadow-xl overflow-hidden w-[300px]">
      <div className="flex w-full bg-gradient-to-t from-sky-700 via-sky-800 to-sky-700">
          <p className="text-white ml-2 text-sm">Link</p>
      </div>
      <form action="/link/api" method="POST">
          <div className="flex flex-col items-start bg-white rounded-md p-4">
              <input id="linkinput" type="text" name="url" className="w-full p-1 shadow-sm rounded-md mb-2 focus:outline-sky-800 border-2 border-solid border-gray-400 transition-all" placeholder="https://example.com" required/>
              <button className="flex p-1 pl-2 pr-2 text-center text-white rounded-md shadow-sm bg-sky-800 active:bg-sky-600 hover:bg-sky-700 transition-all">Shortify link</button>
          </div>
      </form>
    </div>
  );
}
