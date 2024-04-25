<template>
  <el-dialog
    title="创建签名证书"
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
      @submit.native.prevent
    >
      <el-row>
        <el-col :span="24">
          <el-form-item label="证书名称" prop="dn">
            <el-input
              v-model="formData.dn"
              type="text"
              placeholder="请输入证书名称"
              clearable
              @keyup.enter.native="submit"
            />
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
import { domain } from "@/utils/validate";

export default {
  data() {
    return {
      dialogVisible: false,
      formData: {
        dn: ""
      },
      rules: {
        dn: domain
      }
    };
  },
  methods: {
    show() {
      this.dialogVisible = true;
    },
    handleClose() {
      this.dialogVisible = false;
      this.$refs.dataForm.resetFields();
    },
    submit() {
      this.$refs.dataForm.validate(async valid => {
        if (!valid) return false;
        this.dialogVisible = false;
        this.$emit("handleLoading", true);
        try {
          let res = await this.$api.sign.CreateSign(this.formData);
          this.$emit("handleLoading", false);
          this.$message["success"]("创建成功");
        } catch (err) {
          this.$emit("handleLoading", false);
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
