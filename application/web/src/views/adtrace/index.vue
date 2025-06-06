<template>
    <div class="offer-trace-container">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="报价" name="offer"></el-tab-pane>
        <el-tab-pane label="合同" name="contract"></el-tab-pane>
        <el-tab-pane label="账户余额" name="balance"></el-tab-pane>
        <el-tab-pane label="管理员操作" name="adminAction"></el-tab-pane>
      </el-tabs>
      <div v-if="activeTab === 'offer'">
        <div style="margin-bottom: 30px;margin-top: 15px; font-weight: bold; font-size: 40px">
          系统所有报价信息查询
        </div>
        <!-- 搜索栏 -->
        <el-input
          v-model="input1"
          placeholder="请输入报价码查询"
          style="width: 300px;margin-right: 15px;"
        />
        <el-button type="primary" plain @click="handleOfferQuery">查询</el-button>
        <el-button type="success" plain @click="handleAllOffers">获取所有报价信息</el-button>
  
        <!-- 筛选条件 -->
        <el-checkbox style="margin-left: 10px;" v-model="filters1.purchase" label="购电" />
        <el-checkbox v-model="filters1.sale" label="售电" />
  
        <!-- 表格 -->
        <el-table :data="currentPageOffers" style="width: 100%">
          <el-table-column label="报价码" prop="offerId" />
          <el-table-column label="单价/元" prop="offerSnapshot.price" />
          <el-table-column label="数量/kwh" prop="offerSnapshot.quantity" />
          <el-table-column
            label="类型"
            prop="offerSnapshot.isSeller"
            :formatter="formatIsSeller"
          />
          <el-table-column label="创建时间" prop="offerSnapshot.timestamp" />
          <el-table-column label="更新时间" prop="offerSnapshot.updatedTime" />
          <el-table-column label="操作备注" prop="action" />
        </el-table>
  
        <!-- 分页组件 -->
        <el-pagination
          @size-change="handleOfferSizeChange"
          @current-change="handleOfferCurrentChange"
          :current-page="offerCurrentPage"
          :page-sizes="[10, 20, 30, 50]"
          :page-size="offerPageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="offerTotalCount"
        />
      </div>
      <div v-if="activeTab === 'contract'">
        <div style="margin-bottom: 30px;margin-top: 15px; font-weight: bold; font-size: 40px">
          系统所有合同信息查询
        </div>
        <!-- 搜索栏 -->
        <el-input
          v-model="input2"
          placeholder="请输入合同码查询"
          style="width: 300px;margin-right: 15px;"
        />
        <el-button type="primary" plain @click="handleContractQuery">查询</el-button>
        <el-button type="success" plain @click="handleAllContracts">获取所有合同信息</el-button>
  
        <!-- 筛选条件 -->
        <el-checkbox style="margin-left: 10px;" v-model="filters2.purchase" label="购电" />
        <el-checkbox v-model="filters2.sale" label="售电" />
  
        <!-- 表格 -->
        <el-table :data="currentPageContracts" style="width: 100%">
          <el-table-column label="合同编号" prop="contractId" />
          <el-table-column label="售电方" prop="sellerName" />
          <el-table-column label="购电方" prop="buyerName" />
          <el-table-column label="电力单价" prop="price" />
          <el-table-column label="交易电量" prop="quantity" />
          <el-table-column label="交易时间" prop="timestamp" />
        </el-table>
  
        <!-- 分页组件 -->
        <el-pagination
          @size-change="handleContractSizeChange"
          @current-change="handleContractCurrentChange"
          :current-page="contractCurrentPage"
          :page-sizes="[10, 20, 30, 50]"
          :page-size="contractPageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="contractTotalCount"
        />
      </div>
      <div v-if="activeTab === 'balance'">
        <div style="margin-bottom: 30px;margin-top: 15px; font-weight: bold; font-size: 40px">
          管理员账户余额信息查询
        </div>
        <!-- 筛选条件 -->
        <el-checkbox style="margin-left: 10px;" v-model="filters3.out" label="出账" />
        <el-checkbox v-model="filters3.in" label="入账" />
  
        <!-- 表格 -->
        <el-table :data="currentPageBalances" style="width: 100%">
          <el-table-column label="时间" prop="timestamp" />
          <el-table-column label="变化量/元" prop="amount" />
          <el-table-column label="余额/元" prop="rest" />
          <el-table-column label="变动原因" prop="reason" />
          <el-table-column label="来源/去向" prop="userName" />
          <el-table-column label="流向" prop="isIncome" :formatter="formatIsIncome"/>
        </el-table>
  
        <!-- 分页组件 -->
        <el-pagination
          @size-change="handleBalanceSizeChange"
          @current-change="handleBalanceCurrentChange"
          :current-page="balanceCurrentPage"
          :page-sizes="[10, 20, 30, 50]"
          :page-size="balancePageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="balanceTotalCount"
        />
      </div>
      <div v-if="activeTab === 'adminAction'">
        <div style="margin-bottom: 30px;margin-top: 15px; font-weight: bold; font-size: 40px">
          管理员操作溯源信息查询
        </div>
        <!-- 表格 -->
        <el-table :data="currentPageAdminActions" style="width: 100%">
          <el-table-column label="操作" prop="action" />
          <el-table-column label="时间" prop="timestamp" />
          <el-table-column label="详情" prop="details" />
        </el-table>
  
        <!-- 分页组件 -->
        <el-pagination
          @size-change="handleAdminActionSizeChange"
          @current-change="handleAdminActionCurrentChange"
          :current-page="adminActionCurrentPage"
          :page-sizes="[10, 20, 30, 50]"
          :page-size="adminActionPageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="adminActionTotalCount"
        />
      </div>
    </div>
  </template>
  
  <script>
  import { mapGetters } from 'vuex';
  import { getAdminActionHistory, getAdminMoneyHistory, adminGetAllOffer, getAllContract } from '@/api/trace';
  import { getName } from '@/api/user';
  
  export default {
    data() {
      return {
        activeTab: 'offer', // 默认选中的标签
        input1: '',
        input2: '',
        filters1: { purchase: false, sale: false },
        filters2: { purchase: false, sale: false },
        filters3: { out: false, in: false },
        offerhistorydata: [], // 原始全量报价数据
        contractsdata: [], // 原始全量合同数据
        balancehistorydata: [], // 原始全量账户余额数据
        adminActiondata: [], // 原始全量管理员操作数据
        currentPageOffers: [], // 当前页报价数据
        currentPageContracts: [], // 当前页合同数据
        currentPageBalances: [], // 当前页账户余额数据
        currentPageAdminActions: [], // 当前页管理员操作数据
        offerCurrentPage: 1, // 报价当前页码
        offerPageSize: 10, // 报价每页显示数量
        offerTotalCount: 0, // 报价总记录数
        contractCurrentPage: 1, // 合同当前页码
        contractPageSize: 10, // 合同每页显示数量
        contractTotalCount: 0, // 合同总记录数
        balanceCurrentPage: 1, // 账户余额当前页码
        balancePageSize: 10, // 账户余额每页显示数量
        balanceTotalCount: 0, // 账户余额总记录数
        adminActionCurrentPage: 1, // 管理员操作当前页码
        adminActionPageSize: 10, // 管理员操作每页显示数量
        adminActionTotalCount: 0, // 管理员操作总记录数
        loading: false,
        isQuerying1: false,
        isQuerying2: false,
      };
    },
    computed: {
      ...mapGetters([
        'name',
        'userType'
      ]),
      filteredOffers() {
        return this.offerhistorydata.filter(item => {
          if (!this.filters1.purchase && !this.filters1.sale) return true;
          if (this.filters1.purchase && !item.offerSnapshot.isSeller) return true;
          if (this.filters1.sale && item.offerSnapshot.isSeller) return true;
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
        return this.balancehistorydata.filter(item => {
          if (!this.filters3.out && !this.filters3.in) return true;
          if (this.filters3.out && item.amount < 0) return true;
          if (this.filters3.in && item.amount > 0) return true;
          return false;
        });
      },
    },
    created() {
      // 初始化获取第一页报价数据
      this.handleAllOffers();
    },
    methods: {
      // 报价分页数据处理
      handleOfferPagination() {
        const start = (this.offerCurrentPage - 1) * this.offerPageSize;
        const end = this.offerCurrentPage * this.offerPageSize;
        this.currentPageOffers = this.filteredOffers.slice(start, end);
        this.offerTotalCount = this.filteredOffers.length;
      },
  
      // 报价页码变化处理
      handleOfferCurrentChange(page) {
        this.offerCurrentPage = page;
        this.handleOfferPagination();
      },
  
      // 报价每页数量变化处理
      handleOfferSizeChange(size) {
        this.offerPageSize = size;
        this.offerCurrentPage = 1; // 切换每页数量后重置到第一页
        this.handleOfferPagination();
      },
  
      // 报价搜索查询
      handleOfferQuery() {
        this.isQuerying1 = true;
        if (this.input1) {
          const offer = this.offerhistorydata.filter(item => item.offerId === this.input1);
          if (offer.length > 0) {
            this.offerhistorydata = offer;
            this.$message.success('查找成功');
          } else {
            this.$message.warning('未找到对应的报价信息');
          }
        } else {
          this.handleAllOffers(); // 如果搜索框为空，获取所有数据
        }
        this.handleOfferPagination();
      },
  
      // 获取所有报价数据（带分页参数）
      handleAllOffers() {
        this.loading = true;
        adminGetAllOffer().then(res => {
          if (res.code === 200) {
            this.offerhistorydata = JSON.parse(res.data);
            this.offerTotalCount = this.offerhistorydata.length;
            this.handleOfferPagination();
            this.$message.success('获取成功');
          } else {
            this.$message.error('获取数据失败');
          }
          this.loading = false;
        }).catch(err => {
          console.error(err);
          this.$message.error('申请失败');
          this.loading = false;
        });
      },
  
      // 合同分页数据处理
      handleContractPagination() {
        const start = (this.contractCurrentPage - 1) * this.contractPageSize;
        const end = this.contractCurrentPage * this.contractPageSize;
        this.currentPageContracts = this.filteredContracts.slice(start, end);
        this.contractTotalCount = this.filteredContracts.length;
      },
  
      // 合同页码变化处理
      handleContractCurrentChange(page) {
        this.contractCurrentPage = page;
        this.handleContractPagination();
      },
  
      // 合同每页数量变化处理
      handleContractSizeChange(size) {
        this.contractPageSize = size;
        this.contractCurrentPage = 1; // 切换每页数量后重置到第一页
        this.handleContractPagination();
      },
  
      // 合同搜索查询
      handleContractQuery() {
        this.isQuerying2 = true;
        if (this.input2) {
          const contract = this.contractsdata.filter(item => item.contractId === this.input2);
          if (contract.length > 0) {
            this.contractsdata = contract;
            this.$message.success('查找成功');
          } else {
            this.$message.warning('未找到对应的合同信息');
          }
        } else {
          this.handleAllContracts(); // 如果搜索框为空，获取所有数据
        }
        this.handleContractPagination();
      },
  
      // 获取所有合同数据（带分页参数）
      handleAllContracts() {
        this.loading = true;
        getAllContract().then(res => {
          if (res.code === 200) {
            this.contractsdata = JSON.parse(res.data);
            this.contractTotalCount = this.contractsdata.length;
            this.handleContractPagination();
            this.$message.success('获取成功');
          } else {
            this.$message.error('获取数据失败');
          }
          this.loading = false;
        }).catch(err => {
          console.error(err);
          this.$message.error('申请失败');
          this.loading = false;
        });
      },
  
      // 账户余额分页数据处理
      handleBalancePagination() {
        const start = (this.balanceCurrentPage - 1) * this.balancePageSize;
        const end = this.balanceCurrentPage * this.balancePageSize;
        this.currentPageBalances = this.filteredBalances.slice(start, end);
        this.balanceTotalCount = this.filteredBalances.length;
      },
  
      // 账户余额页码变化处理
      handleBalanceCurrentChange(page) {
        this.balanceCurrentPage = page;
        this.handleBalancePagination();
      },
  
      // 账户余额每页数量变化处理
      handleBalanceSizeChange(size) {
        this.balancePageSize = size;
        this.balanceCurrentPage = 1; // 切换每页数量后重置到第一页
        this.handleBalancePagination();
      },
  
      // 获取所有账户余额数据（带分页参数）
      handleAllBalances() {
        this.loading = true;
        getAdminMoneyHistory().then(res => {
          if (res.code === 200) {
            this.balancehistorydata = JSON.parse(res.data);
            this.balanceTotalCount = this.balancehistorydata.length;
            this.handleBalancePagination();
            this.$message.success('获取成功');
          } else {
            this.$message.error('获取数据失败');
          }
          this.loading = false;
        }).catch(err => {
          console.error(err);
          this.$message.error('申请失败');
          this.loading = false;
        });
      },
  
      // 管理员操作分页数据处理
      handleAdminActionPagination() {
        const start = (this.adminActionCurrentPage - 1) * this.adminActionPageSize;
        const end = this.adminActionCurrentPage * this.adminActionPageSize;
        this.currentPageAdminActions = this.adminActiondata.slice(start, end);
        this.adminActionTotalCount = this.adminActiondata.length;
      },
  
      // 管理员操作页码变化处理
      handleAdminActionCurrentChange(page) {
        this.adminActionCurrentPage = page;
        this.handleAdminActionPagination();
      },
  
      // 管理员操作每页数量变化处理
      handleAdminActionSizeChange(size) {
        this.adminActionPageSize = size;
        this.adminActionCurrentPage = 1; // 切换每页数量后重置到第一页
        this.handleAdminActionPagination();
      },
  
      // 获取所有管理员操作数据（带分页参数）
      handleAllAdminActions() {
        this.loading = true;
        getAdminActionHistory().then(res => {
          if (res.code === 200) {
            this.adminActiondata = JSON.parse(res.data);
            this.adminActionTotalCount = this.adminActiondata.length;
            this.handleAdminActionPagination();
            this.$message.success('获取成功');
          } else {
            this.$message.error('获取数据失败');
          }
          this.loading = false;
        }).catch(err => {
          console.error(err);
          this.$message.error('申请失败');
          this.loading = false;
        });
      },
  
      // 格式化报价类型
      formatIsSeller: (row, column, cellValue) => cellValue ? '售电' : '购电',
      // 格式化账户余额流向
      formatIsIncome: (row, column, cellValue) => cellValue ? '入账' : '出账',
      
      // 标签切换处理
      handleTabChange(tab) {
        if (tab === 'offer' && this.offerhistorydata.length === 0) {
          this.handleAllOffers();
        } else if (tab === 'contract' && this.contractsdata.length === 0) {
          this.handleAllContracts();
        } else if (tab === 'balance' && this.balancehistorydata.length === 0) {
          this.handleAllBalances();
        } else if (tab === 'adminAction' && this.adminActiondata.length === 0) {
          this.handleAllAdminActions();
        }
      },
    },
    watch: {
      activeTab: {
        handler: 'handleTabChange',
        immediate: true
      },
      filters1: {
        deep: true,
        handler: 'handleOfferPagination'
      },
      filters2: {
        deep: true,
        handler: 'handleContractPagination'
      },
      filters3: {
        deep: true,
        handler: 'handleBalancePagination'
      }
    }
  };
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