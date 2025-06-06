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
        电力交易报价溯源信息查询
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
        电力交易合同信息查询
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
        用户账户余额溯源信息查询
      </div>
      <!-- 筛选条件 -->
      <el-checkbox style="margin-left: 10px;" v-model="filters3.out" label="出账" />
      <el-checkbox v-model="filters3.in" label="入账" />

      <!-- 表格 -->
      <el-table :data="currentPageBalances" style="width: 100%">
        <el-table-column label="时间" prop="timestamp" />
        <el-table-column label="变化量" prop="amount" />
        <el-table-column label="余额" prop="rest" />
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
import { getOfferHistory, getUserContracts, getBalanceHistory, getAdminActionHistory } from '@/api/trace';
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
      offerFiltersChanged: false, // 报价筛选标志
      contractFiltersChanged: false, // 合同筛选标志
      balanceFiltersChanged: false, // 账户余额筛选标志
      adminActionFiltersChanged: false, // 管理员操作筛选标志
    };
  },
  watch: {
    filters1: {
      deep: true,
      handler() {
        this.offerFiltersChanged = true;
        this.handleOfferPagination(); // 重新计算报价分页数据
      },
    },
    filters2: {
      deep: true,
      handler() {
        this.contractFiltersChanged = true;
        this.handleContractPagination(); // 重新计算合同分页数据
      },
    },
    filters3: {
      deep: true,
      handler() {
        this.balanceFiltersChanged = true;
        this.handleBalancePagination(); // 重新计算账户余额分页数据
      },
    },
  },
  computed: {
    filteredOffers() {
      let filtered = this.offerhistorydata;
      if (this.filters1.purchase || this.filters1.sale) {
        filtered = filtered.filter(item => {
          if (this.filters1.purchase && !item.isSeller) return true;
          if (this.filters1.sale && item.isSeller) return true;
          return false;
        });
      }
      return filtered;
    },
    filteredContracts() {
      return this.contractsdata.filter(contract => {
        const isPurchase = contract.buyerId === this.name;
        const isSale = contract.sellerId === this.name;

        if (this.filters2.purchase && isPurchase) return true;
        if (this.filters2.sale && isSale) return true;

        return !this.filters2.purchase && !this.filters2.sale;
      });
    },
    filteredBalances() {
      return this.balancehistorydata.filter(item => {
        // 如果两个选项都未勾选，显示所有数据
        if (!this.filters3.out && !this.filters3.in) return true;
        if (this.filters3.out && item.amount < 0) return true;
        if (this.filters3.in && item.amount > 0) return true;
        // 其他情况不显示
        return false;
      });
    },
  },
  created() {
    // 初始化获取第一页报价数据
    this.handleAllOffers();
    this.handleAllContracts();
    this.handleAllBalances();
    this.handleAllAdminActions();
  },
  methods: {
    // 报价分页数据处理
    handleOfferPagination() {
      const start = (this.offerCurrentPage - 1) * this.offerPageSize;
      const end = this.offerCurrentPage * this.offerPageSize;
      this.currentPageOffers = this.filteredOffers.slice(start, end);
      this.offerTotalCount = this.filteredOffers.length;
      this.offerFiltersChanged = false; // 重置标志
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
      // 执行查询逻辑（如根据报价码过滤全量数据）
      if (this.input1) {
        this.offerhistorydata = this.offerhistorydata.filter(item => 
          item.offerId === this.input1
        );
      }
      this.handleOfferPagination(); // 重新计算分页
    },

    // 获取所有报价数据（带分页参数）
    handleAllOffers() {
      getOfferHistory().then(res => {
        if (res.code === 200) {
          this.offerhistorydata = JSON.parse(res.data); // 假设接口返回分页数据结构
          this.offerTotalCount = this.offerhistorydata.length;// 总记录数
          this.handleOfferPagination(); // 渲染当前页数据
        } else {
          this.$message.error('获取数据失败');
        }
      } ).catch (err => {
        console.error(err);
        this.$message.error('申请失败');
      });
    },

    // 合同分页数据处理
    handleContractPagination() {
      const start = (this.contractCurrentPage - 1) * this.contractPageSize;
      const end = this.contractCurrentPage * this.contractPageSize;
      this.currentPageContracts = this.filteredContracts.slice(start, end);
      this.contractTotalCount = this.filteredContracts.length;
      this.contractFiltersChanged = false; // 重置标志
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
      // 执行查询逻辑（如根据合同码过滤全量数据）
      if (this.input2) {
        this.contractsdata = this.contractsdata.filter(item => 
          item.contractId === this.input2
        );
      }
      this.handleContractPagination(); // 重新计算分页
    },

    // 获取所有合同数据（带分页参数）
    handleAllContracts() {
      getUserContracts().then(res => {
        if (res.code === 200) {
          this.contractsdata = JSON.parse(res.data);
          console.log(this.contractsdata);
          this.contractTotalCount = this.contractsdata.length;
          this.handleContractPagination();
        } else {
          this.$message.error('获取数据失败');
        }
      }).catch(err => {
        console.error(err);
        this.$message.error('申请失败');
      });
    },

    // 账户余额分页数据处理
    handleBalancePagination() {
      const start = (this.balanceCurrentPage - 1) * this.balancePageSize;
      const end = this.balanceCurrentPage * this.balancePageSize;
      this.currentPageBalances = this.filteredBalances.slice(start, end);
      this.balanceTotalCount = this.filteredBalances.length;
      this.balanceFiltersChanged = false; // 重置标志
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
      getBalanceHistory().then(res => {
        if (res.code === 200) {
          this.balancehistorydata = JSON.parse(res.data);
          console.log(this.balancehistorydata);
          this.balanceTotalCount = this.balancehistorydata.length;
          this.handleBalancePagination();
        } else {
          this.$message.error('获取数据失败');
        }
      }).catch(err => {
        console.error(err);
        this.$message.error('申请失败');
      });
    },

    // 管理员操作分页数据处理
    handleAdminActionPagination() {
      const start = (this.adminActionCurrentPage - 1) * this.adminActionPageSize;
      const end = this.adminActionCurrentPage * this.adminActionPageSize;
      this.currentPageAdminActions = this.adminActiondata.slice(start, end);
      this.adminActionTotalCount = this.adminActiondata.length;
      this.adminActionFiltersChanged = false; // 重置标志
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
      getAdminActionHistory().then(res => {
        if (res.code === 200) {
          this.adminActiondata = JSON.parse(res.data);
          console.log(this.adminActiondata);
          this.adminActionTotalCount = this.adminActiondata.length;
          this.handleAdminActionPagination();
        } else {
          this.$message.error('获取数据失败');
        }
      }).catch(err => {
        console.error(err);
        this.$message.error('申请失败');
      });
    },

    // 格式化报价类型
    formatIsSeller: (row, column, cellValue) => cellValue ? '售电' : '购电',
    // 格式化账户余额流向
    formatIsIncome: (row, column, cellValue) => cellValue ? '入账' : '出账',
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