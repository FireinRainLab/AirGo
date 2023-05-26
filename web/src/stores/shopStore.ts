import { defineStore } from "pinia";
//api
import { useShopApi } from "/@/api/shop/index";
import { ElMessage } from "element-plus";
const shopApi = useShopApi()


export const useShopStore = defineStore("shopStore", {
    state: () => ({
        //商品管理页面数据
        goodsManageData: {
            isShowDialog: false,
            type: '',
            title: '',
            total: 0,
            //loading: true,
            pageNum: 1,
            pageSize: 10,
            //当前编辑商品
            currentGoods: {
                //  id: 0, 不能覆盖
                subject: "",
                total_amount: 0,
                product_code: "",
                total_bandwidth: 0,
                expiration_date: 0,
                checked_nodes:[0], //套餐编辑时选中的节点
                nodes:[],
            } as Goods,
        },

        //全部商品
        goodsList: [] as Goods[],
        //商店页面参数
        tableData: {
            isShowPurchaseDialog: false,
            isShowQRDialog: false,
            isShowSubmitOrderDialog:false,
            //二维码支付链接
            QRcode: null,
            QRtext: '',
            //当前支付商品
            currentGoods: {
                id: 0,
                subject: "",
                total_amount: 0,
                product_code: "",
                total_bandwidth: 0,
                expiration_date: 0,
            } as Goods,
            //当前支付商品创建订单
            currentOrder:{
                id:0,
                out_trade_no:'',
                goods_id:0,
                subject:'',
                price:'',
                pay_type:'alipay',
                trade_no:'',
                buyer_logon_id:'',
                trade_status:'',
                total_amount:'',
               // status:'',
                qr_code:'',
            } as Order,
        }
    }),
    actions: {
        //获取全部订阅商品
        async getAllGoods() {
            const res = await shopApi.getAllGoodsApi()
            // console.log("获取全部订阅商品:",res.data)
            this.goodsList = res.data
            this.goodsManageData.total = this.goodsList.length
        },
        //添加商品
        async newGoods() {
            const res = await shopApi.newGoodsApi(this.goodsManageData.currentGoods)
            if (res.code === 0) {
                ElMessage.success(res.msg)
            } else {
                ElMessage.error(res.msg)
            }

        },
        //修改商品
        async updateGoods() {
            const res = await shopApi.updateGoodsApi(this.goodsManageData.currentGoods)
            if (res.code === 0) {
                ElMessage.success(res.msg)
            } else {
                ElMessage.error(res.msg)
            }
        },
        //删除商品
        async deleteGoods(goods: object) {
            const res = await shopApi.deleteGoodsApi(this.goodsManageData.currentGoods)
            if (res.code === 0) {
                ElMessage.success(res.msg)
            } else {
                ElMessage.error(res.msg)
            }

        },
    }
})