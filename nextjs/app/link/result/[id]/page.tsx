'use server';

import QRCode from 'qrcode'

// import { Description, Dialog, DialogPanel, DialogTitle } from '@headlessui/react'
// import { useState } from 'react'

// function QrCodePopup(qrSrc: string) {
//   let [isOpen, setIsOpen] = useState(false)

//   return (
//     <>
//       <button onClick={() => setIsOpen(true)}>Open dialog</button>
//       <Dialog open={isOpen} onClose={() => setIsOpen(false)} className="relative z-50">
//        <div id="qrbox" className="absolute z-10 flex w-full h-full justify-center items-center bg-[#00000088] invisible cursor-pointer">
//         <div className="flex flex-col items-end p-4 bg-white rounded-xl shadow-lg cursor-default">
//             <div className="w-full flex flex-row justify-between items-center">
//                 <p>QR Code</p>
//                 <div className="w-[32px] h-[32px] rounded-[16px] bg-white text-xl text-bold text-center cursor-pointer hover:bg-gray-200">x</div>
//             </div>
//             <img width="256" height="256" src={qrSrc} alt="QRCode"/>
//         </div>
//       </div>
//       </Dialog>
//     </>
//   )
// }

export default async function Page(
  {params}: {params: Promise<{ id: string }>}
) {
  const { id } = await params;
  const linkUrl = `http:localhost:3000/link/${id}`;
  const qrSrc = await QRCode.toDataURL(linkUrl)

  return (
    <div className="flex flex-col bg-white rounded-md shadow-lg overflow-hidden">
        <div className="flex w-full bg-gradient-to-t from-sky-700 via-sky-800 to-sky-700">
            <p className="text-white ml-2 text-sm">Result</p>
        </div>
        <div className="flex flex-row items-stretch p-4">
            <p 
              id="linkbox" 
              className="flex flex-row bg-gray-100 rounded-md border-2 p-2 text-center border-solid border-gray-200 hover:border-sky-800 cursor-pointer"
            >
              { linkUrl }
            </p>
            <img 
              className="rounded-md ml-2 p-1 border-2 border-solid border-gray-200 hover:border-sky-800 w-[48px] h-[48px] cursor-pointer" 
              src={ qrSrc } 
              alt="QRCode"
              // onClick={}
            />
        </div>
    </div>
  );
}