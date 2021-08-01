export interface Image {
    id: string,
    owner: string,
    title: string,
    description: string,
    created_at: string,
    updated_at: string,
    original_name: string,
    mime_type: string,
    state: ImageState,
    metadata: {
        [key: string]: string
    },
    url: string,
}

export enum ImageState {
    CREATED = "created",
    QUEUED = "queued",
    IN_PROGRESS = "in_progress",
    DONE = "done",
    ERROR = "error"
}
