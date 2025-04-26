import Form from "next/form";
import { notFound, redirect, RedirectType } from "next/navigation";

export default function Page() {
  return (
    <div className="flex flex-col bg-white rounded-md overflow-hidden shadow-xl w-[350px]">
      <div className="flex w-full bg-gradient-to-t from-sky-700 via-sky-800 to-sky-700">
          <p className="text-white ml-2 text-sm">Registration</p>
      </div>
      <Form action="/user/create" formMethod="POST" formEncType="multipart/form-data">
          <div className="flex flex-col items-center bg-white rounded-md p-4">
              <input type="email" name="email" placeholder="user@example.com" className="rounded-md pl-1 mb-2 border-2 border-solid border-gray-400" required/>
              <input type="password" name="password" placeholder="password" className="rounded-md pl-1 mb-2 border-2 border-solid border-gray-400" required/>
              <button id="loginbutton" className="flex p-1 pl-2 pr-2 text-center text-white rounded-md shadow-sm bg-sky-800 active:bg-sky-600 hover:bg-sky-700 transition-all">Create account</button>
          </div>
      </Form>
  </div>
  );
}