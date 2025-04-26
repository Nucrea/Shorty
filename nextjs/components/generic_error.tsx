function GenericError(status: number, text: string) {
    return (
        <div className="flex flex-col items-center">
            <p className="text-white font-bold text-9xl text-center">{status}</p>
            <br/>
            <p className="text-white font-bold text-4xl text-center">{text}</p>
        </div>
    );
}

export default GenericError;