<template>
  <div class="uplink-container">
    <div style="margin-bottom: 30px; font-weight: bold; font-size: 40px">
      用户基本信息
    </div>
    <div class="second-row">
      <div class="user-info">
        <div style="color:#909399;margin-bottom: 30px">
          <table style="width: 100%; border-collapse: collapse; text-align: left;line-height: 30px;">
            <tr>
              <td style=" padding-right: 10px;">当前用户：</td>
              <td>{{ name }}</td>
            </tr>
            <tr>
              <td style=" padding-right: 10px;">用户身份：</td>
              <td>{{ userType }}</td>
            </tr>
            <tr>
              <td style=" padding-right: 10px;">购电资质：</td>
              <td>
                {{ userdata.isBuyer }}
                <el-button v-if="userdata.isBuyer === '未通过'||userdata.isBuyer === '未申请'" type="text" @click="applyForQualification('buy')">
                  申请
                </el-button>
              </td>
              <!-- <td>
                {{ purchaseQualification }}
                <el-button v-if="purchaseQualification === '未通过'" type="text" @click="applyForQualification('buy')">
                  申请
                </el-button>
              </td> -->
            </tr>
            <tr>
              <td style=" padding-right: 10px;">售电资质：</td>
              <td>
                {{ userdata.isSeller }}
                <el-button v-if="userdata.isSeller === '未通过'||userdata.isSeller === '未申请'" type="text" @click="applyForQualification('buy')">
                  申请
                </el-button>
              </td>
              <!-- <td>
                {{ saleQualification }}
                <el-button v-if="saleQualification === '未通过'" type="text" @click="applyForQualification('sell')">
                  申请
                </el-button>
              </td> -->
            </tr>
            <tr>
              <td style=" padding-right: 10px;">账户余额：</td>
              <td>{{ userdata.balance }}/元</td>
            </tr>
          </table>
        </div>
      </div>
      <div class="form-container">
        <div>
          <el-form ref="form" :model="userdata" label-width="80px" size="mini" style="">
  

            <div v-show="userType == '普通用户'">
              <el-form-item label="报价:" style="width: 300px" label-width="120px">
                <el-input v-model="offerdata.price" />
              </el-form-item>
              <el-form-item label="数量:" style="width: 300px" label-width="120px">
                <el-input v-model="offerdata.quantity" />
              </el-form-item>
              <el-form-item label="类型:" style="width: 300px" label-width="120px">
                <el-select v-model="offerdata.isSeller" placeholder="请选择报价类型">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
              </el-form-item>
        
            </div>

          </el-form>
          <span slot="footer" style="color: gray;" class="dialog-footer">
            <el-button v-show="(userdata.isSeller && offerdata.isSeller) || (userdata.isBuyer && !(offerdata.isSeller))"
              type="primary" plain style="margin-left: 220px;" @click="submittracedata()">提 交</el-button>
          </span>
          <span v-show="!((userdata.isSeller && offerdata.isSeller) || (userdata.isBuyer && !(offerdata.isSeller)))"
            slot="footer" style="color: gray;" class="dialog-footer">
            没有权限提交报价！请先完成资质审核!
          </span>
        </div>
      </div>
    </div>
    <div class="search-container">
      <el-input v-model="input" placeholder="请输入报价码查询" style="width: 300px;margin-right: 15px;" />
      <el-button type="primary" plain @click="OfferInfo"> 查询 </el-button>
      <el-button type="success" plain @click="AllOfferInfo"> 获取所有报价信息 </el-button>
      <!-- 勾选购电、售电类型的筛选选项 -->
      <div v-if="!isQuerying" style="margin-top: 15px;">
        <el-checkbox v-model="filters.purchase" label="购电" />
        <el-checkbox v-model="filters.sale" label="售电" />
      </div>
      <el-table :data="filteredOfferData" style="width: 100%">
        <el-table-column label="报价码" prop="offerID" />
        <el-table-column label="类型" prop="isSeller" :formatter="formatIsSeller" />
        <el-table-column label="报价" prop="price" />
        <el-table-column label="数量" prop="quantity" />
        <el-table-column label="时间" prop="timestamp" />
        <el-table-column label="状态" prop="status" />
        <el-table-column label="撮合次数" prop="round" />
      </el-table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { uplink, userApproveAs, userGetAllOffer } from '@/api/trace'
import { getUserInfo } from '@/api/user'

export default {
  name: 'Uplink',
  data() {
    return {
      userdata: {
        userID: '',
        balance: 0,
        isSeller: '',
        isBuyer: '',
        offers: [],
        contracts: [],
        balanceHistory: [],
        offerHistory: [],
        offerDone: [],
        creditRating:0,
        tradeCount:0,
      },
      offerdata: {
        offerID: '',
        userID: '',
        price: 0,
        quantity: 0,
        deposit: 0,
        isSeller: false,
        timestamp: '',
        status: '',
        round: 0
      },
      isQuerying: false,
      AllOffers: [],
      input: '',
      loading: false,
      options: [{
        value: false,
        label: '购电'
      }, {
        value: true,
        label: '售电'
      }],
      filters: {
        purchase: false, // 购电筛选选项
        sale: false // 售电筛选选项
      }
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'userType'
    ]),
    // 计算属性，根据筛选条件动态筛选表格数据
    filteredOfferData() {
      if (this.isQuerying) {
        return [this.offerdata];
      }
      return this.getAllOffers().filter(item => {
        // 如果两个选项都未勾选，显示所有数据
        if (!this.filters.purchase && !this.filters.sale) return true;

        // 如果勾选了购电，筛选类型为“购电”的数据
        if (this.filters.purchase && !item.isSeller) return true;

        // 如果勾选了售电，筛选类型为“售电”的数据
        if (this.filters.sale && item.isSeller) return true;

        // 其他情况不显示
        return false;
      });
    },
    // purchaseQualification() {
    //   if (this.userdata.isBuyer) {
    //     return '已通过审核';
    //   } else if (!this.userdata.isBuyer && !this.userdata.approveUserAsBuyer) {
    //     return '未通过';
    //   } else if (!this.userdata.isBuyer && this.userdata.approveUserAsBuyer) {
    //     return '正在申请';
    //   }
    //   return '';
    // },
    // saleQualification() {
    //   if (this.userdata.isSeller) {
    //     return '已通过审核';
    //   } else if (!this.userdata.isSeller && !this.userdata.approveUserAsSeller) {
    //     return '未通过';
    //   } else if (!this.userdata.isSeller && this.userdata.approveUserAsSeller) {
    //     return '正在申请';
    //   }
    //   return '';
    // },

    // getOffer(){
    //   //合并userdata的Offers和OfferDone数组，并遍历查找并返回特定的报价码对应的offer
    //   //return append(this.userdata.Offers,this.userdata.OfferDone...)
    // }
  },
  created() {
    getUserInfo().then(res => {
      console.log("created")
      this.userdata = JSON.parse(res.data)
      console.log(this.userdata)
    })
  },
  methods: {
    getAllOffers() {
      return [...this.userdata.offers, ...this.userdata.offerDone];
    },
    getOffer() {
      return (offerId) => {
        return this.getAllOffers().find(item => item.offerID === offerId);
      };
    },
    submittracedata() {
      console.log(this.userdata)
      const loading = this.$loading({
        lock: true,
        text: '数据上链中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      var formData = new FormData()
      //formData.append('traceability_code', this.tracedata.traceability_code)
      // 根据不同的用户给arg1、arg2、arg3..赋值,

      formData.append('arg1', this.offerdata.price)
      formData.append('arg2', this.offerdata.quantity)
      formData.append('arg3', this.offerdata.isSeller)

      uplink(formData).then(res => {
        if (res.code === 200) {
          loading.close()
          this.$message({
            message: '上链成功，交易ID：' + res.txid + '\n报价码：' + res.offerId,
            type: 'success'
          })
        } else {
          loading.close()
          this.$message({
            message: '上链失败',
            type: 'error'
          })
        }
      }).catch(err => {
        loading.close()
        console.log(err)
      })
    },
    applyForQualification(status) {
      // 调用申请资质的接口
      console.log('申请资质');
      const loading = this.$loading({
        lock: true,
        text: '申请中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      // 示例：调用后端接口申请资质
      const formData = new FormData();
      formData.append('status', status);
      userApproveAs(formData).then(res => {
        if (res.code === 200) {
          loading.close()
          this.$message.success('申请成功');
          if (status === 'sell') {
            this.userdata.isSeller = '申请中';
          } else {
            this.userdata.isBuyer = '申请中';
          } // 更新状态
        } else {
          loading.close()
          this.$message.error('申请失败');
        }
      }).catch(err => {
        console.error(err);
        loading.close()
        this.$message.error('申请失败');
      });
    },

    OfferInfo() {
      // 根据输入的报价码查询数据
      console.log('查询报价码:', this.input);
      this.isQuerying = true;
      // 这里调用getOffer()
      const offer = this.getOffer(this.input);
      if (offer) {
        console.log('查询结果:', offer);
        // 可以在这里更新界面显示查询结果
        this.$message.success('查找成功');
        this.offerdata = offer;
      } else {
        this.$message.warning('未找到对应的报价信息');
      }
    },
    AllOfferInfo() {
      // 获取所有报价信息
      console.log('获取所有报价信息');
      this.isQuerying = false;
      const loading = this.$loading({
        lock: true,
        text: '获取报价中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      // 这里可以调用后端接口获取数据
      userGetAllOffer().then(res => {
        if (res.code === 200) {
          loading.close()
          this.AllOffers = JSON.parse(res.data);
          this.$message.success('申请成功');
        } else {
          loading.close()
          this.$message.error('申请失败');
        }
      }).catch(err => {
        console.error(err);
        loading.close()
        this.$message.error('申请失败');
      });
    },
    formatIsSeller(row, column, cellValue) {
      return cellValue ? '售电' : '购电';
    }
  }
}

</script>

<style lang="scss" scoped>
.uplink {
  &-container {
    display: block;
    // justify-content: space-between;
    // align-items: flex-start;
    margin: 30px;
  }

  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}

.second-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 30px;
}

.user-info {
  width: 48%; // 调整宽度以适应布局
}

.form-container {
  width: 48%; // 调整宽度以适应布局
}

.search-container {
  margin-top: 30px;
}
</style>
