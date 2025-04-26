import { redirect, RedirectType } from "next/navigation";

export async function POST(request: Request) {
    console.log(request.url)
    return redirect('https://example.com', RedirectType.push);
}