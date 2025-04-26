import api from "@/api/api";
import { notFound, redirect, RedirectType } from "next/navigation";

export default async function Page(
  {params}: {params: Promise<{ id: string }>}
) {
  const { id } = await params;
  const result = await api.GetLink(id);
  if (result == null) {
    return notFound();
  }

  return redirect(result.url, RedirectType.push);
}