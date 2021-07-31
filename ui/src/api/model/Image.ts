export interface Image {
    id: string,
    owner: string,
    title: string,
    description: string,
    created_at: string,
    updated_at: string,
    original_name: string,
    mime_type: string,
    state: string,
    metadata: {
        [key: string]: string
    },
    url: string,
}
