import GenericError from "@/components/generic_error";

export default function UnauthorizedPage() {
    return GenericError(401, "Unauthorized");
}