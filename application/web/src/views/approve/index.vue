<template>
  <div>
    <div class="uplink-container">
      <div class="countdown-container" style="margin-bottom: 20px; padding: 10px; background-color: #f5f7fa; border-radius: 4px;">
        <div style="font-size: 14px; color: #909399; margin-top: 10px;">
          <span>下一次撮合时间：</span>
          <span style="color: #303133; font-weight: 500;">{{ nextMatchTimeFormatted }}</span>
          <!-- <span style="margin-left: 10px;">（撮合频率：{{ matchFrequency }}分钟/次）</span> -->
        </div>
      </div>
      <div style="margin-bottom: 30px; font-weight: bold; font-size: 40px">
        管理员基本信息
      </div>

      <div class="second-row">
        <div class="user-info-container">
          <!-- 管理员信息表格 -->
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
                  <td style=" padding-right: 10px;">账户余额：</td>
                  <td>{{ admindata.balance.toFixed(2)  }}元</td>
                </tr>
              </table>
            </div>
          </div>
          <!-- 系统配置信息表格 -->
          <div class="form-container">
            <!-- <div style="margin-bottom: 10px; font-weight: bold; font-size: 30px">
              系统配置信息
            </div> -->
            <div style="color:#909399;margin-bottom: 30px">
              <table style="width: 100%; border-collapse: collapse; text-align: left;line-height: 30px;">
                <tr>
                  <td style="padding-right: 10px;">撮合时间：</td>
                  <td>
                    <span v-if="!isEditing.MatchFrequency" style="display: inline-block; width: 100px;">{{ adminconfig.matchFrequency }}分钟/次</span>
                    <input v-else v-model="tempMatchFrequency" type="text" style="width: 100px; box-sizing: border-box;" />
                  </td>
                  <td>
                    <button v-if="!isEditing.MatchFrequency" @click="startEdit('MatchFrequency')">修改</button>
                    <button v-else @click="submitEdit('MatchFrequency')">提交</button>
                  </td>
                </tr>
                <tr>
                  <td style="padding-right: 10px;">保证金率：</td>
                  <td>
                    <span v-if="!isEditing.DepositRate" style="display: inline-block; width: 100px;">{{ adminconfig.depositRate }}</span>
                    <input v-else v-model="tempDepositRate" type="number" style="width: 100px; box-sizing: border-box;"/>
                  </td>
                  <td>
                    <button v-if="!isEditing.DepositRate" @click="startEdit('DepositRate')">修改</button>
                    <button v-else @click="submitEdit('DepositRate')">提交</button>
                  </td>
                </tr>
                <tr>
                  <td style="padding-right: 10px;">手续费率：</td>
                  <td>
                    <span v-if="!isEditing.FeeRate" style="display: inline-block; width: 100px;">{{ adminconfig.feeRate }}</span>
                    <input v-else v-model="tempFeeRate" type="number" style="width: 100px; box-sizing: border-box;"/>
                  </td>
                  <td>
                    <button v-if="!isEditing.FeeRate" @click="startEdit('FeeRate')">修改</button>
                    <button v-else @click="submitEdit('FeeRate')">提交</button>
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="search-container">
      <el-input v-model="input" placeholder="请输入用户名查询资质" style="width: 300px;margin-right: 15px;" />
      <el-button type="primary" plain @click="ApproveList"> 查询 </el-button>
      <el-button type="success" plain @click="AllApproveList"> 获取所有用户资质 </el-button>
      <div style="margin-top: 15px;">
        <el-checkbox v-model="filters.purchase" label="购电申请" @change="handleFilterChange" />
        <el-checkbox v-model="filters.sale" label="售电申请" @change="handleFilterChange" />
      </div>
      <el-table :data="displayedApplicationData" style="width: 100%">
        <el-table-column label="用户ID" prop="userId" />
        <el-table-column label="申请类型" prop="applyType" :formatter="formatIsSeller"/>
        <el-table-column label="申请时间" prop="applyTime"/>
        <el-table-column label="详情">
          <template #default="scope">
            <el-button type="info" @click="viewDetails(scope.row)">查看</el-button>
          </template>
        </el-table-column>
        <el-table-column label="审核操作">
          <template #default="scope">
            <div v-if="scope.row.auditStatus === '审核中'">
              <el-radio-group v-model="scope.row.auditResult">
                <el-radio :label="true">通过</el-radio>
                <el-radio :label="false">拒绝</el-radio>
              </el-radio-group>
              <el-button type="primary" @click="submitApproval(scope.row)">确定</el-button>
            </div>
            <div v-else>
              {{ scope.row.auditStatus }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="审核时间" prop="auditTime"/>

      </el-table>
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[10, 20, 30]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
      </el-pagination>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getConfig, adminModify, approveUserAs } from '@/api/trace'
import { getAdminInfo, getName } from '@/api/user'
export default {
  name: 'Approve',
  data() {
    return {
      admindata: {
        adminId: '',
        balance: 0,
        balanceHistory: [],
        adminActionHistory: [],
        applications:[],
        //sellList: [],
        //buyList: [],
        contracts: [],
      },
      isEditing: {
        MatchFrequency: false,
        DepositRate: false,
        FeeRate: false,
      },
      tempMatchFrequency: '',
      tempDepositRate: '',
      tempFeeRate: 0,
      adminconfig: {
        matchFrequency: 0,
        depositRate: 0,
        feeRate: 0,
      },
      isQuerying: false,
      loading: false,
      input: '',
      filters: {
        purchase: false,
        sale: false
      },
      uinfoData: [],
      displayedApplicationData: [],
      currentPage: 1, // 当前页码
      pageSize: 10,   // 每页显示的数量
      total: 0 ,       // 总记录数
      matchFrequency: 0,        // 撮合频率（分钟）
      nextMatchTime: null,          // 下次撮合时间（Date对象）
      nextMatchTimeFormatted: '',   // 格式化后的下次撮合时间
    };
  },
  computed: {
    ...mapGetters([
        'name',
        'userType'
      ]),
    filteredApplicationData() {
      console.log("filter applications", this.filters, this.admindata.applications);
      
      if (this.isQuerying) {
        return this.uinfoData;
      }
      
      return this.admindata.applications.filter(item => {
        // 如果两个选项都未勾选，显示所有数据
        if (!this.filters.purchase && !this.filters.sale) return true;
        
        // 根据申请类型和筛选条件进行匹配
        if (this.filters.purchase && item.applyType === 'buy') return true;
        if (this.filters.sale && item.applyType === 'sell') return true;
        
        return false;
      });
    },
  },
  created() {
      this.fetchAdminInfo();
  },
  
  methods: {
    viewDetails(row) {
      console.log('查看详细信息:', row);
      // 这里可以添加弹出模态框显示详细信息的逻辑
    },
    fetchAdminInfo() {
      getAdminInfo().then(res => {
          console.log("admin info loaded", res);
          this.admindata = JSON.parse(res.data);
          // 初始化申请数据的审核结果
          this.admindata.applications.forEach(app => {
            app.auditResult = app.auditStatus === '审核通过';
          });
          this.getConfig();
      }).catch(err => {
          console.error('获取管理员信息失败:', err);
      });
    },
    
    getConfig() {
      getConfig().then(res => {
        this.adminconfig = JSON.parse(res.data);
        console.log("admin config loaded", this.adminconfig);
        this.matchFrequency = this.adminconfig.matchFrequency;
        this.updateNextMatchTime();
        this.AllApproveList();
      }).catch(err => {
          console.error('获取配置信息失败:', err);
      });
    },
    
    updateNextMatchTime() {
      const now = new Date();
      const minutes = now.getMinutes();
      const remainder = minutes % this.matchFrequency;
      const minutesToNextMatch = remainder === 0 ? this.matchFrequency : this.matchFrequency - remainder;
      
      // 计算下次撮合的时间
      this.nextMatchTime = new Date(now);
      this.nextMatchTime.setMinutes(now.getMinutes() + minutesToNextMatch);
      this.nextMatchTime.setSeconds(0); // 秒数设为0，精确到分钟
      
      // 格式化下次撮合时间
      this.formatNextMatchTime();
    },

    // 添加时间格式化方法
    formatNextMatchTime() {
      if (!this.nextMatchTime) return;
      
      const year = this.nextMatchTime.getFullYear();
      const month = String(this.nextMatchTime.getMonth() + 1).padStart(2, '0');
      const day = String(this.nextMatchTime.getDate()).padStart(2, '0');
      const hours = String(this.nextMatchTime.getHours()).padStart(2, '0');
      const minutes = String(this.nextMatchTime.getMinutes()).padStart(2, '0');
      const seconds = String(this.nextMatchTime.getSeconds()).padStart(2, '0');
      
      this.nextMatchTimeFormatted = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },

    formatIsSeller(row, column, cellValue) {
      if(cellValue === 'buy'){
        return '购电';
      }
      if(cellValue === 'sell'){
        return '售电';
      }
      return cellValue;
    },
    
    handleFilterChange() {
      console.log("Filters changed:", this.filters);
      this.currentPage = 1; // 重置页码为1
      this.total = this.filteredApplicationData.length;
      this.handleCurrentChange(1);
    },
    
    submitEdit(field) {
      console.log("Submitting edit for field:", field, this.adminconfig);
      const loading = this.$loading({
        lock: true,
        text: '数据修改中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      
      var formData = new FormData();
      if (field === 'MatchFrequency') {
        this.adminconfig.matchFrequency = this.tempMatchFrequency;
        formData.append('name', 'MatchFrequency');
        formData.append('newConfig', this.tempMatchFrequency);
      } else if (field === 'DepositRate') {
        this.adminconfig.depositRate = this.tempDepositRate;
        formData.append('name', 'DepositRate');
        formData.append('newConfig', this.tempDepositRate);
      } else if (field === 'FeeRate') {
        this.adminconfig.feeRate = this.tempFeeRate;
        formData.append('name', 'FeeRate');
        formData.append('newConfig', this.tempFeeRate);
      }
      
      this.isEditing[field] = false;
      
      adminModify(formData).then(res => {
        if (res.code === 200) {
          loading.close();
          this.$message({
            message: '上链成功，交易ID：' + res.txid + '\n',
            type: 'success'
          });
          this.getConfig(); // 重新加载配置
        } else {
          loading.close();
          this.$message({
            message: '上链失败',
            type: 'error'
          });
        }
      }).catch(err => {
        loading.close();
        console.log(err);
      });
    },
    
    startEdit(field) {
      this.isEditing[field] = true;
      if (field === 'MatchFrequency') {
        this.tempMatchFrequency = this.adminconfig.matchFrequency;
      } else if (field === 'DepositRate') {
        this.tempDepositRate = this.adminconfig.depositRate;
      } else if (field === 'FeeRate') {
        this.tempFeeRate = this.adminconfig.feeRate;
      }
    },
    
    submitApproval(row) {
      console.log("Submitting approval:", row);
      const loading = this.$loading({
        lock: true,
        text: '资质审核-数据上链中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      
      var formData = new FormData();
      formData.append('arg1', row.applicationId);
      formData.append('arg2', row.auditResult);

      approveUserAs(formData).then(res => {
        if (res.code === 200) {
          loading.close();
          this.$message({
            message: '上链成功，交易ID：' + res.txid + '\n',
            type: 'success'
          });
          // 更新本地数据状态
          row.auditStatus = row.auditResult ? '审核通过' : '审核拒绝';
          row.auditTime = new Date().toISOString();
          // 重新加载管理员信息
          this.fetchAdminInfo();
        } else {
          loading.close();
          this.$message({
            message: '上链失败',
            type: 'error'
          });
        }
      }).catch(err => {
        loading.close();
        console.log(err);
      });
    },
    
    getUinfo() {
      return (uname) => {
        return this.admindata.applications.filter(item => item.userId === uname);
      };
    },
    
    ApproveList() {
      // 查询逻辑
      console.log('查询用户名:', this.input);
      const finduinfo = this.getUinfo();
      const uinfo = finduinfo(this.input);
      if (uinfo.length > 0) {
        console.log('查询结果:', uinfo);
        // 可以在这里更新界面显示查询结果
        this.$message.success('查找成功');
        this.uinfoData = uinfo;
        this.isQuerying = true;
        this.total = uinfo.length;
        this.handleCurrentChange(1); // 重置页码为1
      } else {
        this.$message.warning('未找到对应的用户申请信息');
      }
    },
    
    AllApproveList() {
      // 获取所有用户资质逻辑
      this.filters.purchase = false;
      this.filters.sale = false;
      console.log("Loading all applications:", this.admindata.applications);
      this.isQuerying = false;
      this.input = '';
      this.total = this.filteredApplicationData.length;
      this.handleCurrentChange(1); // 重置页码为1
    },
    
    handleSizeChange(newSize) {
      this.pageSize = newSize;
      this.handleCurrentChange(this.currentPage);
    },
    
    handleCurrentChange(newPage) {
      this.currentPage = newPage;
      const start = (newPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      this.displayedApplicationData = this.filteredApplicationData.slice(start, end);
    }
  }
};
</script>

<style lang="scss" scoped>
* {
  box-sizing: border-box;
}
// 为整个管理员信息和系统配置信息容器设置样式
.uplink-container {
  display: block;
  margin: 30px;
}

// 第二行（包含管理员信息和系统配置信息的行）
.second-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 30px;
}

// 包含管理员信息表格和系统配置信息表格的容器
.user-info-container {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

// 管理员信息表格容器
.user-info {
  width: 48%; 
}

// 系统配置信息表格容器
.form-container {
  width: 48%; 
}

// 为表格设置基本样式，使其布局稳定
table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  line-height: 30px;
}

// 为表格中的单元格设置样式
td {
  padding-right: 10px;
}

// 编辑状态下输入框的样式
input {
  width: 40%; // 根据实际情况调整宽度
  padding: 3px;
  box-sizing: border-box;
}

// 编辑状态下按钮的样式
button {
  padding: 5px 10px;
}
.search-container {
  margin-top: 30px;
  margin-left: 20px;
}
</style>  