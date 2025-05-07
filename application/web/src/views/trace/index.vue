<template>
  <div class="trace-container">
      <div class="button-container">
          <el-button type="primary" style="margin-bottom: 20px; font-size: 40px; font-weight: bold;"
              @click="offerInfo">电力交易报价溯源信息查询</el-button>
          <div v-if="showoffer">
              <el-input v-model="input1" placeholder="请输入报价码查询" style="width: 300px;margin-right: 15px;" />
              <el-button type="primary" plain @click="OfferInfo"> 查询 </el-button>
              <el-button type="success" plain @click="AllOfferInfo"> 获取所有报价信息 </el-button>
              <el-checkbox v-model="filters1.purchase" label="购电" />
              <el-checkbox v-model="filters1.sale" label="售电" />
              <el-table :data="filteredOffers" style="width: 100%">
                  <el-table-column label="报价码" prop="offer.offerID" />
                  <el-table-column label="单价" prop="offer.price" />
                  <el-table-column label="数量" prop="offer.quantity" />
                  <el-table-column label="类型" prop="offer.isSeller" :formatter="formatIsSeller" />
                  <el-table-column label="操作时间" prop="timestamp" />
                  <el-table-column label="操作备注" prop="action" />
              </el-table>
          </div>
          <el-button type="primary" style="margin-bottom: 20px; font-size: 40px; font-weight: bold;"
              @click="contractInfo">电力交易合同信息查询</el-button>
          <div v-if="showcontract">
              <el-input v-model="input2" placeholder="请输入合同码查询" style="width: 300px;margin-right: 15px;" />
              <el-button type="primary" plain @click="querryContract"> 查询 </el-button>
              <el-button type="success" plain @click="fetchContracts"> 获取所有合同信息 </el-button>
              <el-checkbox v-model="filters2.purchase" label="购电" />
              <el-checkbox v-model="filters2.sale" label="售电" />
              <el-table :data="filteredContracts" style="width: 100%">
                  <el-table-column label="合同编号" prop="contractID" />
                  <el-table-column label="售电方" prop="sellerName" />
                  <el-table-column label="购电方" prop="buyerName" />
                  <el-table-column label="电力单价" prop="price" />
                  <el-table-column label="交易电量" prop="quantity" />
                  <el-table-column label="交易时间" prop="timestamp" />
              </el-table>
          </div>
          <el-button type="primary" style="margin-bottom: 20px; font-size: 40px; font-weight: bold;"
              @click="balanceInfo">用户账户余额溯源信息查询</el-button>
          <div v-if="showbalance">
              <el-checkbox v-model="filters3.out" label="出账" />
              <el-checkbox v-model="filters3.in" label="入账" />
              <el-table :data="filteredBalances" style="width: 100%">
                  <el-table-column label="时间" prop="timestamp" />
                  <el-table-column label="变化量" prop="amount" />
                  <el-table-column label="余额" prop="rest" />
                  <el-table-column label="变动原因" prop="reason" />
              </el-table>
          </div>
          <el-button type="primary" style="margin-bottom: 20px; font-size: 40px; font-weight: bold;"
              @click="adminInfo">管理员操作溯源信息查询</el-button>
          <div v-if="showadmin">
              <el-table :data="adminActiondata" style="width: 100%">
                  <el-table-column label="操作" prop="action" />
                  <el-table-column label="时间" prop="timestamp" />
                  <el-table-column label="详情" prop="details" />
              </el-table>
          </div>
      </div>

  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { getOfferHistory, getUserContracts, getBalanceHistory, getAdminActionHistory } from '@/api/trace';
import { getName } from '@/api/user';

export default {
  data() {
      return {
          showoffer: false,
          showcontract: false,
          showbalance: false,
          showadmin: false,
          offerhistorydata: [],
          contractsdata: [],
          balancehistorydata: [],
          adminActiondata: [],
          loading: false,
          input1: '',
          input2: '',
          isQuerying1: false,
          isQuerying2: false,
          filters1: {
              purchase: false,
              sale: false
          },
          filters2: {
              purchase: false,
              sale: false
          },
          filters3: {
              out: false,
              in: false
          }
      }
  },
  computed: {
      ...mapGetters([
          'name',
          'userType'
      ]),
      filteredOffers() {
          return this.offerhistorydata.filter(item => {
              // 如果两个选项都未勾选，显示所有数据
              if (!this.filters1.purchase && !this.filters1.sale) return true;
              // 如果勾选了购电，筛选类型为“购电”的数据
              if (this.filters1.purchase && !item.isSeller) return true;
              // 如果勾选了售电，筛选类型为“售电”的数据
              if (this.filters1.sale && item.isSeller) return true;
              // 其他情况不显示
              return false;
          });
      },
      filteredContracts() {
          return this.contractsdata.filter(contract => {
              const isPurchase = contract.buyerName === this.name;
              const isSale = contract.sellerName === this.name;

              if (this.filters2.purchase && isPurchase) return true;
              if (this.filters2.sale && isSale) return true;

              return !this.filters2.purchase && !this.filters2.sale;
          });
      },
      filteredBalances() {
          return this.balancehistorydata().filter(item => {
              // 如果两个选项都未勾选，显示所有数据
              if (!this.filters3.out && !this.filters3.in) return true;
              if (this.filters3.out && item.amount < 0) return true;
              if (this.filters3.in && item.amount > 0) return true;
              // 其他情况不显示
              return false;
          });
      },
      getOffer() {
          return (offerId) => {
              return this.AllOfferInfo().filter(item => item.offerID === offerId);
          };
      },
      getContract() {
          return (Id) => {
              return this.fetchContracts().filter(item => item.contractID === Id);
          };
      },
  },
  // created() {

  // },
  methods: {
      offerInfo() {
          this.showoffer = !this.showoffer
          // 电力交易报价溯源信息查询逻辑
          if (this.showoffer && this.offerhistorydata.length === 0) {
              this.AllOfferInfo();
          }
          console.log('显示电力交易报价溯源信息');
      },
      contractInfo() {
          // 电力交易合同信息查询逻辑
          this.showcontract = !this.showcontract
          if (this.showcontract && this.contractsdata.length === 0) {
              this.fetchContracts();
          }
          console.log('显示电力交易合同信息');
      },
      balanceInfo() {
          // 用户账户余额溯源信息查询逻辑
          this.showbalance = !this.showbalance
          if (this.showbalance && this.balancehistorydata.length === 0) {
              this.AllBalanceInfo();
          }
          console.log('显示用户账户余额溯源信息');
      },
      adminInfo() {
          // 管理员操作溯源信息查询逻辑
          this.showadmin = !this.showadmin
          if (this.showadmin && this.adminActiondata.length === 0) {
              //getAdminActionHistory();
              this.getAction();
          }
          console.log('显示管理员操作溯源信息');
      },
      formatIsSeller(row, column, cellValue) {
          return cellValue ? '售电' : '购电';
      },
      fetchContracts() {
          getUserContracts().then(res => {
              const contracts = JSON.parse(res.data);
              Promise.all(contracts.map(contract => {
                  return Promise.all([
                      this.getUserInfo(contract.sellerID),
                      this.getUserInfo(contract.buyerID)
                  ]).then(([sellerName, buyerName]) => {
                      return {
                          ...contract,
                          sellerName: sellerName,
                          buyerName: buyerName
                      };
                  });
              })).then(updatedContracts => {
                  this.contractsdata = updatedContracts;
              }).catch(err => {
                  console.error('获取用户信息失败:', err);
                  this.$message.error('获取用户信息失败');
              });
          }).catch(err => {
              console.error('获取合同信息失败:', err);
              this.$message.error('获取合同信息失败');
          });
      },
      getAction() {
          console.log('获取管理员操作信息');
          getAdminActionHistory().then(res => {
              if (res.code === 200) {
                  this.adminActiondata = JSON.parse(res.data);
                  this.$message.success('申请成功');
              } else {
                  this.$message.error('申请失败');
              }
          }).catch(err => {
              console.error(err);
              this.$message.error('申请失败');
          });
      },
      getUserInfo(id) {
          var formData = new FormData()
          formData.append('id', id)
          return getName(formData).then(res => {
              return res.data; // 假设返回的数据中包含用户名
          }).catch(err => {
              console.error('获取用户信息失败:', err);
              return '未知用户'; // 如果获取失败，返回默认值
          });
      },
      querryContract() {
          // 根据输入的报价码查询数据
          console.log('查询报价码:', this.input2);
          this.isQuerying2 = true;
          // 这里调用getOffer()
          const offer = this.getContract(this.input2);
          if (offer.length > 0) {
              console.log('查询结果:', offer);
              // 可以在这里更新界面显示查询结果
              this.$message.success('查找成功');
              this.contractsdata = offer;
          } else {
              this.$message.warning('未找到对应的报价信息');
          }
      },
      OfferInfo() {
          // 根据输入的报价码查询数据
          console.log('查询报价码:', this.input1);
          this.isQuerying1 = true;
          // 这里调用getOffer()
          const offer = this.getOffer(this.input1);
          if (offer.length > 0) {
              console.log('查询结果:', offer);
              // 可以在这里更新界面显示查询结果
              this.$message.success('查找成功');
              this.offerhistorydata = offer;
          } else {
              this.$message.warning('未找到对应的报价信息');
          }
      },
      AllOfferInfo() {
          // 获取所有报价信息
          console.log('获取所有报价信息');
          this.isQuerying1 = false;
          // 这里可以调用后端接口获取数据
          getOfferHistory().then(res => {
              if (res.code === 200) {
                  this.offerhistorydata = JSON.parse(res.data);
                  console.log(this.offerhistorydata)
                  //this.$message.success('申请成功');
              } else {
                  this.$message.error('申请失败');
              }
          }).catch(err => {
              console.error(err);
              this.$message.error('申请失败');
          });
      },
      AllBalanceInfo() {
          console.log('获取所有账户余额信息');
          getBalanceHistory().then(res => {
              if (res.code === 200) {
                  this.balancehistorydata = JSON.parse(res.data);
                  this.$message.success('申请成功');
              } else {
                  this.$message.error('申请失败');
              }
          }).catch(err => {
              console.error(err);
              this.$message.error('申请失败');
          });
      },
  }
}
</script>

<style lang="scss" scoped>
.trace {
  &-container {
      display: block;
      margin: 30px;
  }

  &-text {
      font-size: 30px;
      line-height: 46px;
  }
}

.button-container {
  display: flex;
  flex-direction: column;
  width: fit-content;
}
</style>