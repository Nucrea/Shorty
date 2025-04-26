import GenericError from "@/components/generic_error";

export default function ForbiddenPage() {
    return GenericError(403, "Forbidden");
}