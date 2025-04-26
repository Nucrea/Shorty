interface LinkResult {
    id: string
    url: string
}

interface ImageResult {
    id: string
    name: string
    size: number
    originalUrl: string
    thumbnailUrl: string
}

class Api {
    shortyUrl: string;
    apiKey: string;

    constructor(shortyUrl: string, apiKey: string) {
        this.shortyUrl = shortyUrl
        this.apiKey = apiKey
    }

    async CreateLink(inputUrl: string): Promise<LinkResult> {
        const result = await fetch(
            `${this.shortyUrl}/api/v1/link/create`, {
                method: 'post', 
                body: JSON.stringify({url: inputUrl}) 
            });
        const body = await result.json();

        if (result.status < 200 || result.status >= 300) {
            const { message } = body['error'];
            throw Error(message);
        }
        
        const { id, url } = body['result'];
        return { id: id, url: url };
    }

    async GetLink(linkId: string): Promise<LinkResult | null> {
        const result = await fetch(`${this.shortyUrl}/api/v1/link/${linkId}`);
        const body = await result.json();

        if (result.status < 200 || result.status >= 300) {
            const { message } = body['error'];
            throw Error(message);
        }
        if (result.status == 404) {
            return null
        }
        
        const { id, url } = body['result'];
        return { id, url };
    }

    async UploadImage(file: File): Promise<ImageResult | null> {
        const formData = new FormData();
        formData.append('file', file);

        const result = await fetch(
            `${this.shortyUrl}/api/v1/image/upload`, {
                method: 'post',
                // headers: {
                //     'Content-Type': 'multipart/form-data'
                // },
                body: formData,
            });
        const body = await result.json();

        if (result.status < 200 || result.status >= 300) {
            const { message } = body['error'];
            throw Error(message);
        }
        if (result.status == 404) {
            return null
        }
        
        const { id, name, size, originalUrl, thumbnailUrl } = body['result'];
        return { id, name, size, originalUrl, thumbnailUrl };
    }

    async GetImageInfo(imageId: string): Promise<ImageResult | null> {
        const result = await fetch(`${this.shortyUrl}/api/v1/image/info/${imageId}`);
        const body = await result.json();

        if (result.status < 200 || result.status >= 300) {
            const { message } = body['error'];
            throw Error(message);
        }
        if (result.status == 404) {
            return null
        }
        
        const { id, name, size, originalUrl, thumbnailUrl } = body['result'];
        return { id, name, size, originalUrl, thumbnailUrl };
    }
}

export default new Api(process.env.NEXT_SHORTY_URL!, process.env.NEXT_SHORTY_APIKEY!);