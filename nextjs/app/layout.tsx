import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Shorty",
};

function MenuItem(target: string, text: string) {
  return (
    <a href={target} target="_self"
      className="text-white text-center pl-2 pr-2 h-full hover:bg-sky-700 transition-all"
    >
      {text}
    </a>
  );
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <div className="absolute top-0 flex flex-row justify-center items-center w-full h-[40px] bg-sky-800 p-1 shadow-xl">
          { MenuItem("/link", "Link") }
          <div className="border-r ml-2 mr-2 h-full border-white"></div>
          { MenuItem("/image", "Image") }
          <div className="border-r ml-2 mr-2 h-full border-white"></div>
          { MenuItem("/file", "File") }
        </div>
        <div id="logo" className="absolute top-0 left-2">
          <a href="/" target="_self">
            <div className="flex flex-row items-center">
              <img src="/favicon.ico" className="size-6"/>
              <p className="font-bold text-white text-3xl ml-2">Shorty</p>
            </div>
          </a>
        </div>
        <div className="absolute right-2 flex flex-row justify-center items-center h-[40px] bg-sky-800 p-1 shadow-xl">
          {/* if data.Account == nil {
              <a href="/login" target="_self" className="text-white text-center pl-2 pr-2 h-full hover:bg-sky-700 transition-all">Login</a>
          } else {
              <a href="/account" target="_self" className="text-white text-center text-lg pl-2 pr-2 mr-2 h-full hover:bg-sky-700 transition-all">{ data.Account.Email }</a>
              <form action="/logout" method="POST">
                  <button className="text-white text-center pl-2 pr-2 h-full hover:bg-sky-700 transition-all">Logout</button>
              </form>
          } */}
        </div>
        <div className="w-full min-h-screen flex justify-center items-center bg-sky-400">
          { children }
        </div>
        {/* <div className="flex absolute bottom-2 left-1/2 transform -translate-x-1/2">
            <p className="text-white text-sm">{ renderTimeStr(data) }</p>
        </div> */}
      </body>
    </html>
  );
}
