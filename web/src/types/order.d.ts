declare interface Order {
    created_at: string;
    updated_at: string;
    id: number;
    userID: number;
    user_name:string;
    user: any;

    out_trade_no:string;
    goods_id: number;
    subject: string;
    price:string;
    pay_type:string;
   // status:string;

    qr_code:string;
    trade_no: string;
    buyer_logon_id: string;
    trade_status: string;
    total_amount: string;
    receipt_amount: string;
    buyer_pay_amount: string;
}
declare interface OrdersWithTotal{
    order_list:Order[];
    total:number;
}