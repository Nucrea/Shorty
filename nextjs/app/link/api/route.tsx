import api from "@/api/api";
import { redirect, RedirectType } from "next/navigation";
import { NextRequest } from "next/server";

export async function POST(request: NextRequest) {
    const data = await request.formData();
    const formValueUrl = data.get('url');
    const url = formValueUrl?.toString()!;

    const result = await api.CreateLink(url);

    return redirect(`http://localhost:3001/link/result/${result.id}`, RedirectType.replace);
}