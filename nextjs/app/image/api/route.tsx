import api from "@/api/api";
import { redirect, RedirectType } from "next/navigation";

export async function POST(request: Request) {
    const formData = await request.formData();
    const file = formData.get('image') as File;

    console.log(file.name)
    console.log(file.size)

    const result = await api.UploadImage(file);
    if (result) {
        return redirect(`http://localhost:3001/image/view/${result.id}`, RedirectType.replace);
    }
    return 'error occured';
}