declare interface Isp {
    created_at: string;
    updated_at: string;
    id: number;
    user_id: number;
    mobile: string;
    isp_type: string;
    status: boolean;
    unicom_config: UnicomConfig;
}

declare interface UnicomConfig {
    version: string;
    app_id: string;
    cookie: string;
    mobile: string;
    password: string;
}