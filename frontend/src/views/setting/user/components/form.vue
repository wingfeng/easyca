<template>
  <el-dialog
    :title="opt ? '编辑用户/角色' : '创建用户/角色'"
    :visible.sync="dialogVisible"
    width="50%"
    @close="handleClose"
  >
    <el-form
      ref="dataForm"
      :rules="rules"
      :model="formData"
      label-width="120px"
      label-suffix=":"
    >
      <el-row>
        <el-col :span="24">
          <el-form-item label="用户/角色" prop="user">
            <el-input
              v-model="formData.user"
              type="text"
              placeholder="请输入用户/角色名"
              clearable
              :disabled="opt ? true : false"
            />
          </el-form-item>
        </el-col>
        <!-- <el-col :span="24">
          <el-form-item label="权限">
            <el-select
              v-model="formData.roles"
              placeholder="请选择权限"
              multiple
              style="width: 100%"
              clearable
            >
              <el-option
                v-for="(item, idx) in roleList"
                :key="idx"
                :label="item.Label"
                :value="item.Url"
              ></el-option>
            </el-select>
          </el-form-item>
        </el-col> -->
        <el-col :span="24">
          <el-form-item label="权限">
            <el-transfer
              v-model="formData.roles"
              :data="roleList"
              :titles="['权限列表', '拥有权限']"
            ></el-transfer>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取 消</el-button>
      <el-button type="primary" @click="submit">确 定</el-button>
    </span>
  </el-dialog>
</template>

<script>
export default {
  data() {
    return {
      dialogVisible: false,
      opt: 0, // 1: edit option  0: add option
      roleList: [],
      formData: {
        user: "",
        roles: []
      },
      rules: {
        user: [
          {
            required: true,
            message: "请输入用户名",
            trigger: "blur"
          }
        ]
      }
    };
  },
  created() {
    this.getRoleList();
  },
  methods: {
    async show(row) {
      this.dialogVisible = true;
      this.opt = row ? 1 : 0;
      if (row) {
        try {
          let res = await this.$api.user.GetRoleByUser({ user: row.user });
          this.formData = res;
        } catch (error) {
          if (error.data.errInfo) {
            this.$message["warning"](error.data.errInfo);
          } else {
            this.$message["error"](error.message);
          }
        }
      }
    },
    handleClose() {
      this.dialogVisible = false;
      this.formData.roles = [];
      this.$refs.dataForm.resetFields();
    },
    async getRoleList() {
      try {
        let res = await this.$api.user.getRoleList();
        let data = res.list.map(item => {
          return {
            key: item.Url,
            label: item.Label,
            disabled: false
          };
        });
        this.roleList = data;
      } catch (error) {
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
      }
    },
    submit() {
      this.$refs.dataForm.validate(async valid => {
        if (!valid) return false;
        try {
          let res = await this.$api.user.ModifyRole(this.formData);
          this.handleClose();
          this.$emit("handleLoading", false);
          if (this.opt) {
            this.$message["success"]("修改成功");
          } else {
            this.$message["success"]("创建成功");
          }
        } catch (err) {
          if (error.data.errInfo) {
            this.$message["warning"](error.data.errInfo);
          } else {
            this.$message["error"](error.message);
          }
        }
      });
    }
  }
};
</script>

<style lang="less" scoped></style>
