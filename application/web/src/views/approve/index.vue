<template>
    <div>
        <div class="uplink-container">
            <div style="margin-bottom: 30px; font-weight: bold; font-size: 40px">
                管理员基本信息
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
                                <td style=" padding-right: 10px;">账户余额：</td>
                                <td>{{ admindata.balance }}/元</td>
                            </tr>
                        </table>
                    </div>
                    <div class="form-container">
                        <div style="color:#909399;margin-bottom: 30px">
                            <table style="width: 100%; border-collapse: collapse; text-align: left;line-height: 30px;">
                                <tr>
                                    <td style="padding-right: 10px;">撮合时间：</td>
                                    <td>
                                        <span v-if="!isEditing.MatchFrequency">{{ adminconfig.matchFrequency
                                        }}\分钟</span>
                                        <input v-else v-model="tempMatchFrequency" type="text" />
                                    </td>
                                    <td>
                                        <button v-if="!isEditing.MatchFrequency"
                                            @click="startEdit('MatchFrequency')">修改</button>
                                        <button v-else @click="submitEdit('MatchFrequency')">提交</button>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding-right: 10px;">保证金率：</td>
                                    <td>
                                        <span v-if="!isEditing.DepositRate">{{ adminconfig.depositRate }}</span>
                                        <input v-else v-model="tempDepositRate" type="number" />
                                    </td>
                                    <td>
                                        <button v-if="!isEditing.DepositRate"
                                            @click="startEdit('DepositRate')">修改</button>
                                        <button v-else @click="submitEdit('DepositRate')">提交</button>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="padding-right: 10px;">手续费率：</td>
                                    <td>
                                        <span v-if="!isEditing.FeeRate">{{ adminconfig.feeRate }}</span>
                                        <input v-else v-model="tempFeeRate" type="number" />
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
                <el-checkbox v-model="filters.purchase" label="购电申请" />
                <el-checkbox v-model="filters.sale" label="售电申请" />
            </div>
            <el-table :data="displayedOfferData" style="width: 100%">
                <el-table-column label="用户名" prop="username" />
                <el-table-column label="申请类型" prop="IsSeller" :formatter="formatIsSeller" />
                <el-table-column label="审核操作">
                    <template #default="scope">
                        <div v-if="scope.row.reviewStatus === null">
                            <el-checkbox v-model="scope.row.approved" label="通过" />
                            <el-checkbox v-model="!scope.row.approved" label="不通过" />
                            <el-button type="primary" @click="submitApproval(scope.row)">提交审核</el-button>
                        </div>
                        <div v-else>
                            {{ scope.row.reviewStatus === true ? '已审核通过' : '已审核未通过' }}
                        </div>
                    </template>
                </el-table-column>
            </el-table>
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
                adminID: '',
                balance: 0,
                balanceHistory: [],
                adminActionHistory: [],
                sellList: [],
                buyList: [],
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
            displayedOfferData: [],
        };
    },
    computed: {
        ...mapGetters([
            'name',
            'userType'
        ]),

    },
    created() {
        getAdminInfo().then(res => {
            console.log("admin created")
            this.admindata = JSON.parse(res.data)
            console.log(this.admindata)
            this.config()
            this.AllApproveList()
        }).catch(err => {
            console.error('获取管理员信息失败:', err);
        });
    },
    methods: {
        async filteredOfferData() {
            if (this.filters.purchase && this.filters.sale) {
                return [];
            }
            if (this.isQuerying) {
                data = this.uinfodata
            }
            let data = [];
            if (this.filters.purchase) {
                data = this.admindata.buyList.map(item => ({ id: item, IsSeller: false, reviewStatus: null, approved: false }));
            } else if (this.filters.sale) {
                data = this.admindata.sellList.map(item => ({ id: item, IsSeller: true, reviewStatus: null, approved: false }));
            } else {
                data = [...this.admindata.buyList.map(item => ({ id: item, IsSeller: false, reviewStatus: null, approved: false })), ...this.admindata.sellList.map(item => ({ id: item, IsSeller: true, reviewStatus: null, approved: false }))];
            }
            const results = await Promise.all(data.map(async item => {
                var formData = new FormData()
                formData.append('id', item.id)
                const username = await getName(formData);
                //const username=JSON.parse(res.data)
                console.log(username)
                return { ...item, username: username.data };
            }));
            this.displayedOfferData = results;
        },
        config() {
            getConfig().then(res => {
                this.adminconfig = JSON.parse(res.data)
                console.log(this.adminconfig)
            })
        },
        submitEdit(field) {
            console.log(this.adminconfig)
            const loading = this.$loading({
                lock: true,
                text: '数据修改中...',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            })
            var formData = new FormData()
            if (field === 'MatchFrequency') {
                this.adminconfig.matchFrequency = this.tempMatchFrequency;
                formData.append('name', 'MatchFrequency')
                formData.append('newConfig', this.tempMatchFrequency)
            } else if (field === 'DepositRate') {
                this.adminconfig.depositRate = this.tempDepositRate;
                formData.append('name', 'DepositRate')
                formData.append('newConfig', this.tempDepositRate)
            } else if (field === 'FeeRate') {
                this.adminconfig.feeRate = this.tempFeeRate;
                formData.append('name', 'FeeRate')
                formData.append('newConfig', this.tempFeeRate)
            }
            this.isEditing[field] = false;
            adminModify(formData).then(res => {
                if (res.code === 200) {
                    loading.close()
                    this.$message({
                        message: '上链成功，交易ID：' + res.txid + '\n',
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
        formatIsSeller(row) {
            return row.IsSeller ? '售电' : '购电';
        },
        submitApproval(row) {
            console.log(row)
            const loading = this.$loading({
                lock: true,
                text: '资质审核-数据上链中...',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            })
            var formData = new FormData()
            formData.append('arg1', row.id)
            var status = ''
            if (row.IsSeller) {
                status = 'sell'
            } else {
                status = 'buy'
            }
            formData.append('arg2', status)
            formData.append('arg3', row.approved)

            approveUserAs(formData).then(res => {
                if (res.code === 200) {
                    loading.close()
                    this.$message({
                        message: '上链成功，交易ID：' + res.txid + '\n',
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
            row.reviewStatus = row.approved;
        },
        getUinfo() {
            return (uname) => {
                return [...this.admindata.buyList, ...this.admindata.sellList].filter(item => item === uname);
            };
        },
        ApproveList() {
            // 查询逻辑
            console.log('查询用户名:', this.input);
            const uinfo = this.getUinfo(this.input);
            if (uinfo.length > 0) {
                console.log('查询结果:', uinfo);
                // 可以在这里更新界面显示查询结果
                this.$message.success('查找成功');
                this.uinfoData = uinfo;
            } else {
                this.$message.warning('未找到对应的用户申请信息');
            }
        },
        AllApproveList() {
            // 获取所有用户资质逻辑
            this.filters.purchase = false
            this.filters.sale = false
            this.filteredOfferData()
        },
    }
};
</script>

<style lang="scss" scoped>
.uplink {
    &-container {
        display: block;
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