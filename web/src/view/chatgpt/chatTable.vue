<template>
  <div class="gva-table-box">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="连接版" name="first">
        <div v-if="!chatToken">
          <warning-bar title="在资源权限中将此角色的资源权限清空 或者不包含创建者的角色 即可屏蔽此客户资源的显示" />
          <el-input v-model="skObj.sk" class="query-ipt" placeholder="请输入您的ChatGpt SK" clearable />
          <el-button type="primary" @click="save">保存</el-button>
          <div class="secret">
            <el-empty description="请到gpt网站获取您的sk：https://platform.openai.com/account/api-keys" />
          </div>
        </div>
        <div v-else>
          <el-form :model="form" label-width="120px">
            <el-form-item label="数据库类型">
              <el-select v-model="form.dbtype" placeholder="请选择数据库类型" style="width: 115px">
                <el-option
                  v-for="(item, index) in dbTypeArr"
                  :key="index"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="数据库连接">
              <el-input
                v-model="form.url"
                :autosize="{ minRows:1, maxRows: 2 }"
                type="textarea"
                clearable
                placeholder="请输入数据库连接"
                style="width: 400px;"
              />
            </el-form-item>
            <el-form-item label="数据库账号">
              <el-input
                v-model="form.username"
                :autosize="{ minRows:1, maxRows: 2 }"
                type="textarea"
                clearable
                placeholder="请输入数据库账号"
                style="width: 200px;"
              />
            </el-form-item>
            <el-form-item label="数据库密码">
              <el-input
                v-model="form.password"
                :autosize="{ minRows:1, maxRows: 2 }"
                type="password"
                clearable
                placeholder="请输入数据库账号"
                style="width: 200px;"
              />
            </el-form-item>
            <el-form-item label="">
              <el-button type="primary" @click="testConnection">连接数据库</el-button>
            </el-form-item>
            <el-form-item label="查询db名称：">
              <el-select v-model="form.dbname" placeholder="请选择库" style="width: 115px">
                <el-option
                  v-for="(item, index) in dbArr"
                  :key="index"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-if="form.dbtype === 'PostgreSQL'" label="查询schema名称：">
              <el-select v-model="form.schema" placeholder="请选择schema" style="width: 115px">
                <el-option
                  v-for="(item, index) in schemaArr"
                  :key="index"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="查询db描述：">
              <el-input
                v-model="form.chat"
                :autosize="{ minRows: 2, maxRows: 4 }"
                type="textarea"
                clearable
                placeholder="请输入对话"
              />
            </el-form-item>
            <el-form-item label="GPT生成SQL:">
              <el-input
                  v-model="sql"
                  :autosize="{ minRows: 2, maxRows: 4 }"
                  type="textarea"
                  disabled
                  placeholder="此处展示自动生成的sql"
              />
            </el-form-item>
            <el-button type="primary" @click="handleQueryTable">查询</el-button>
          </el-form>
          <div class="tables">
            <el-table
              v-if="tableData.length"
              ref="multipleTable"
              :data="tableData"
              style="width: 100%"
              tooltip-effect="dark"
              height="400px"
            >
              <el-table-column
                v-for="(item, index) in tableData[0]"
                :key="index"
                :prop="index"
                :label="index"
                min-width="200"
                show-overflow-tooltip
              />
            </el-table>
            <p v-else class="text">请在对话框输入你需要AI帮你查询的内容：）</p>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="免连接" name="second">
        <div class="content">
          <p>免连接内容</p>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import { getTableApi,
  createSKApi,
  getSKApi,
  deleteSKApi,
  testConnectApi } from '@/api/chatgpt'
import { getDB as getDBAPI } from '@/api/autoCode'
import { ref, reactive } from 'vue'

const chatToken = ref(null)
const skObj = reactive({
  sk: '',
})
const sql = ref("")
const getSK = async() => {
  const res = await getSKApi()
  chatToken.value = res.data.ok
}

const getDB = async() => {
  //const res = await getDBAPI()
  // if (res.code === 0) {
  //   dbArr.value = res.data.dbs
  // }
}
getSK()
getDB()
const save = async() => {
  const res = await createSKApi(skObj)
  if (res.code === 0) {
    await getSK()
  }
}

const deleteSK = async() => {
  const res = await deleteSKApi()
  if (res.code === 0) {
    await getSK()
  }
}

const form = ref({
  dbtype: '',
  url:'',
  username:'',
  password:'',
  dbname: '',
  schema: '',
  chat: '',
})
const dbArr = ref([])
const schemaArr = ref([])
const dbTypeArr = ref(['MySql', 'PostgreSQL'])
const tableData = ref([])
const activeName = ref('first')

const handleClick = (tab, event) => {
  console.log(tab, event);
}

const handleQueryTable = async() => {
  const res = await getTableApi(form.value)
  if (res.code === 0) {
    tableData.value = res.data.results||[]
  }
  sql.value = res.data.sql
  // 根据后台返回值动态渲染表格
}

const testConnection = async() => {
  const res = await testConnectApi(form.value)
  if (res.code === 0) {
    dbArr.value = res.data.names
    schemaArr.value = res.data.schemas
  }
}
</script>

<style scoped lang="scss">
.secret{
  padding: 30px;
  margin-top: 20px;
  background: #F5F5F5;
  p {
    line-height: 30px;
  }
}
.query-ipt{
  width: 300px;
  margin-right: 30px;
}
.content{
  p {
    font-size: 16px;
    line-height: 20px;
  }
  padding: 10px;
  width: 100%;
  background: #F5F5F5;
  margin-top: 30px;
}
.tables{
  width: 100%;
  margin-top: 30px;
  .text{
    font-size: 18px;
    color: #6B7687;
    text-align: center;
  }
}
</style>
