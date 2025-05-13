<template>
  <div class="uplink-container">
     <!-- 添加消息提示区域 -->
     <div v-if="message" class="notification">
      {{ message }}
    </div>
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
                <el-button v-if="userdata.isSeller === '未通过'||userdata.isSeller === '未申请'" type="text" @click="applyForQualification('sell')">
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
              <td>{{ userdata.balance.toFixed(2) }}/元</td>
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
            <el-button v-show="(userdata.isSeller==='已通过' && offerdata.isSeller) || (userdata.isBuyer==='已通过' && !(offerdata.isSeller))"
              type="primary" plain style="margin-left: 220px;" @click="submitofferdata()">提 交</el-button>
          </span>
          <span v-show="!((userdata.isSeller==='已通过' && offerdata.isSeller) || (userdata.isBuyer==='已通过' && !(offerdata.isSeller)))"
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
        <el-table-column label="报价码" prop="offerId" />
        <el-table-column label="类型" prop="isSeller" :formatter="formatIsSeller" />
        <!-- <el-table-column label="报价/元" prop="price" />
        <el-table-column label="数量" prop="quantity" /> -->
        <el-table-column label="报价/元">
          <template #default="scope">
            <el-input v-if="scope.row.isEditing" v-model="scope.row.price" />
            <span v-else>{{ scope.row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="数量">
          <template #default="scope">
            <el-input v-if="scope.row.isEditing" v-model="scope.row.quantity" />
            <span v-else>{{ scope.row.quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="押金" prop="deposit" />
        <el-table-column label="创建时间" prop="timestamp" />
        <el-table-column label="更新时间" prop="updatedTime" />
        <el-table-column label="状态" prop="status" />
        <el-table-column label="撮合次数" prop="round" />
        <el-table-column label="操作">
          <template #default="scope">
            <span v-if="scope.row.status === '待撮合'">
              <el-button v-if="!scope.row.isEditing" type="primary" size="mini" @click="startEdit(scope.row)">修改</el-button>
              <el-button v-if="!scope.row.isEditing" type="danger" size="mini" @click="cancelOffer(scope.row)">撤销</el-button>
              <el-button v-if="scope.row.isEditing" type="success" size="mini" @click="submitEdit(scope.row)">提交</el-button>
            </span>
            <span v-else>不可操作</span>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { uplink, userApproveAs, userCancel, userGetAllOffer,usermodify } from '@/api/trace'//GetAllOffer
import { getUserInfo } from '@/api/user'

export default {
  name: 'Uplink',
  data() {
    return {
      userdata: {
        userId: '',
        userName: '',
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
        offerId: '',
        userId: '',
        userName: '',
        price: 0,
        quantity: 0,
        deposit: 0,
        isSeller: false,
        timestamp: '',
        updateTime: '',
        status: '',
        round: 0,
        isEditing: false,
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
      },
      // 新增：WebSocket相关变量
      socket: null,
      message: '',
      messageType: 'info' // info, success, warning, error
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
      return this.AllOffers.filter(item => {
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
  },
  created() {
    getUserInfo().then(res => {
      //console.log("created")
      this.userdata = JSON.parse(res.data)
      this.AllOfferInfo()
      console.log(this.userdata)
    }),
    // 新增：初始化WebSocket连接
    this.initWebSocket()
  },
  beforeDestroy() {
    // 新增：组件销毁前关闭WebSocket连接
    this.closeWebSocket()
  },
  methods: {
    getOffer() {
      return (offerId) => {
        return this.AllOffers.find(item => item.offerId === offerId);
      };
    },
    // 新增修改报价方法
    startEdit(row) {
      row.isEditing = true;
    },
    // 新增提交修改方法
    submitEdit(row) {
      const loading = this.$loading({
        lock: true,
        text: '数据上链中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      const formData = new FormData();
      formData.append('offerId', row.offerId);
      formData.append('arg1', row.price);
      formData.append('arg2', row.quantity);

      usermodify(formData).then(res => {
        loading.close()
        if (res.code === 200) {
          this.$message.success('修改成功');
          row.isEditing = false;
          getUserInfo().then(res => {
            //console.log("created")
            this.userdata = JSON.parse(res.data)
            //console.log(this.userdata)
            this.AllOfferInfo()
          })
        } else {
          this.$message.error('修改失败');
        }
      }).catch(err => {
        loading.close()
        console.error(err);
        this.$message.error('修改失败');
      });
    },
    // 新增撤销报价方法
    cancelOffer(row) {
      console.log('撤销报价');
      const loading = this.$loading({
        lock: true,
        text: '数据上链中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      // 这里可以添加撤销报价的逻辑
      const formData = new FormData();
      formData.append('offerId', row.offerId);
      formData.append('userId', row.userId);

      userCancel(formData).then(res => {
        loading.close()
        if (res.code === 200) {
          this.$message.success('撤销成功');
          row.isEditing = false;
          getUserInfo().then(res => {
            //console.log("created")
            this.userdata = JSON.parse(res.data)
            //console.log(this.userdata)
            this.AllOfferInfo()
          })
        } else {
          this.$message.error('撤销失败');
        }
      }).catch(err => {
        loading.close()
        console.error(err);
        this.$message.error('撤销失败');
      });
    },
    submitofferdata() {
      //console.log(this.userdata)
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
          // // 使用后端返回的 offer 数据
          // const newOffer = JSON.parse(res.data);
          // console.log('newOffer:');
          // console.log(newOffer);
          // this.userdata.offers.push(newOffer);
          getUserInfo().then(res => {
            //console.log("created")
            this.userdata = JSON.parse(res.data)
            console.log(this.userdata)
            this.AllOfferInfo()
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
      // 这里调用getOffer()
      //const offer = this.getOffer(this.input);
      const findOfferFn = this.getOffer(); // 先获取返回的箭头函数
      const offer = findOfferFn(this.input); // 调用箭头函数并传入报价码获取实际报价
      if (offer) {
        console.log('查询结果');
        // // 动态添加 isEditing 属性
        offer.isEditing = false;
        console.log(offer);

        // 可以在这里更新界面显示查询结果
        this.$message.success('查找成功');
        this.offerdata = offer;
        this.isQuerying = true;
      } else {
        this.$message.warning('未找到对应的报价信息');
      }
    },
    AllOfferInfo() {
      // 获取所有报价信息
      console.log('获取所有报价信息');
      this.isQuerying = false;
      this.input=''
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
          this.AllOffers = JSON.parse(res.data).map(item => ({
            ...item,
            isEditing: false // 确保每一行都有 isEditing 属性
          }));
          //this.$message.success('申请成功');
        } else {
          loading.close()
          //this.$message.error('申请失败');
        }
      }).catch(err => {
        console.error(err);
        loading.close()
        //this.$message.error('申请失败');
      });
    },
    formatIsSeller(row, column, cellValue) {
      return cellValue ? '售电' : '购电';
    },
    // 新增：初始化WebSocket连接
    initWebSocket() {
      // 检查浏览器是否支持WebSocket
      if (!window.WebSocket) {
        this.showMessage('您的浏览器不支持WebSocket，无法接收实时消息', 'error')
        return
      }
      
      // 创建WebSocket连接（根据实际域名和端口修改）
      this.socket = new WebSocket('ws://localhost:8080/ws')
      
      // 连接成功回调
      this.socket.onopen = () => {
        console.log('WebSocket连接已建立')
        this.showMessage('已连接到服务器，接收实时报价信息', 'info')
      }
      
      // 接收到消息回调
      // this.socket.onmessage = (event) => {
      //   console.log('收到消息:', event.data)
        
      //   // 显示消息
      //   if (event.data.includes("开始")) {
      //     this.showMessage(event.data, 'info')
      //   } else if (event.data.includes("完成")) {
      //     this.showMessage(event.data, 'success')
      //   } else {
      //     this.showMessage(event.data, 'info')
      //   }
        
      //   // 如果是撮合完成消息，刷新报价列表
      //   if (event.data.includes("撮合报价完成")) {
      //     this.AllOfferInfo()
      //   }
      // }

      this.socket.onmessage = (event) => {
    console.log('收到消息:', event.data)

    // 解析 JSON 数据
    try {
        const messageData = JSON.parse(event.data)
        console.log('解析后的消息:', messageData)

        // 提取需要的信息
        const messageType = messageData.type || 'info'
        const message = messageData.message || '未知消息'

        // 显示消息
        this.showMessage(message, messageType)

        // 如果是撮合完成消息，刷新报价列表
        if (message.includes("撮合报价完成")) {
            this.AllOfferInfo()
        }
    } catch (error) {
        console.error('解析消息失败:', error)
        this.showMessage('收到未知格式的消息', 'warning')
    }
}
      
      // 连接关闭回调
      this.socket.onclose = (event) => {
        console.log('WebSocket连接已关闭:', event)
        this.showMessage('与服务器的连接已断开，正在尝试重连...', 'warning')
        
        // 尝试重连
        setTimeout(() => {
          this.initWebSocket()
        }, 3000)
      }
      
      // 连接错误回调
      this.socket.onerror = (error) => {
        console.error('WebSocket错误:', error)
        this.showMessage('与服务器的连接发生错误', 'error')
      }
    },
    
    // 新增：关闭WebSocket连接
    closeWebSocket() {
      if (this.socket) {
        this.socket.close()
        this.socket = null
      }
    },
    
    // 新增：显示消息提示
    showMessage(content, type = 'info') {
      this.message = content
      this.messageType = type
      
      // 3秒后自动隐藏消息
      clearTimeout(this.messageTimer)
      this.messageTimer = setTimeout(() => {
        this.message = ''
      }, 5000)
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

/* 新增：消息提示样式 */
.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  color: white;
  transition: all 0.3s ease;
  
  &.info {
    background-color: #409eff;
  }
  
  &.success {
    background-color: #67c23a;
  }
  
  &.warning {
    background-color: #e6a23c;
  }
  
  &.error {
    background-color: #f56c6c;
  }
}
</style>
