export interface StageContext {
    api: API,
    ses: SES
}

export interface API {
    stageName: string
}

export interface SES {
    fromEmail: string
}
