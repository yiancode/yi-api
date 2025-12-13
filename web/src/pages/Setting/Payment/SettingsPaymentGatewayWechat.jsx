/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useEffect, useState, useRef } from 'react';
import {
  Banner,
  Button,
  Form,
  Row,
  Col,
  Typography,
  Spin,
} from '@douyinfe/semi-ui';
const { Text } = Typography;
import {
  API,
  removeTrailingSlash,
  showError,
  showSuccess,
} from '../../../helpers';
import { useTranslation } from 'react-i18next';

export default function SettingsPaymentGatewayWechat(props) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    WechatAppId: '',
    WechatMchId: '',
    WechatApiV2Key: '',
    WechatMinTopUp: 1,
  });
  const [originInputs, setOriginInputs] = useState({});
  const formApiRef = useRef(null);

  useEffect(() => {
    if (props.options && formApiRef.current) {
      const currentInputs = {
        WechatAppId: props.options.WechatAppId || '',
        WechatMchId: props.options.WechatMchId || '',
        WechatApiV2Key: props.options.WechatApiV2Key || '',
        WechatMinTopUp:
          props.options.WechatMinTopUp !== undefined
            ? parseFloat(props.options.WechatMinTopUp)
            : 1,
      };
      setInputs(currentInputs);
      setOriginInputs({ ...currentInputs });
      formApiRef.current.setValues(currentInputs);
    }
  }, [props.options]);

  const handleFormChange = (values) => {
    setInputs(values);
  };

  const submitWechatSetting = async () => {
    if (props.options.ServerAddress === '') {
      showError(t('请先填写服务器地址'));
      return;
    }

    setLoading(true);
    try {
      const options = [];

      if (inputs.WechatAppId && inputs.WechatAppId !== '') {
        options.push({ key: 'WechatAppId', value: inputs.WechatAppId });
      }
      if (inputs.WechatMchId && inputs.WechatMchId !== '') {
        options.push({ key: 'WechatMchId', value: inputs.WechatMchId });
      }
      if (inputs.WechatApiV2Key && inputs.WechatApiV2Key !== '') {
        options.push({ key: 'WechatApiV2Key', value: inputs.WechatApiV2Key });
      }
      if (
        inputs.WechatMinTopUp !== undefined &&
        inputs.WechatMinTopUp !== null
      ) {
        options.push({
          key: 'WechatMinTopUp',
          value: inputs.WechatMinTopUp.toString(),
        });
      }

      // 发送请求
      const requestQueue = options.map((opt) =>
        API.put('/api/option/', {
          key: opt.key,
          value: opt.value,
        }),
      );

      const results = await Promise.all(requestQueue);

      // 检查所有请求是否成功
      const errorResults = results.filter((res) => !res.data.success);
      if (errorResults.length > 0) {
        errorResults.forEach((res) => {
          showError(res.data.message);
        });
      } else {
        showSuccess(t('更新成功'));
        // 更新本地存储的原始值
        setOriginInputs({ ...inputs });
        props.refresh?.();
      }
    } catch (error) {
      showError(t('更新失败'));
    }
    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      <Form
        initValues={inputs}
        onValueChange={handleFormChange}
        getFormApi={(api) => (formApiRef.current = api)}
      >
        <Form.Section text={t('微信支付设置')}>
          <Text>
            微信支付商户号、API密钥等设置请
            <a
              href='https://pay.weixin.qq.com/'
              target='_blank'
              rel='noreferrer'
            >
              点击此处
            </a>
            登录微信支付商户平台进行配置。
            <br />
          </Text>
          <Banner
            type='info'
            description={`支付回调地址填：${props.options.ServerAddress ? removeTrailingSlash(props.options.ServerAddress) : t('网站地址')}/api/wechat/notify`}
          />
          <Banner
            type='warning'
            description={`请确保服务器外网可访问，否则无法接收微信支付回调`}
          />
          <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='WechatAppId'
                label={t('微信 AppID（可选）')}
                placeholder={t('微信公众号 AppID，用于 JSAPI 支付')}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='WechatMchId'
                label={t('商户号')}
                placeholder={t('微信支付商户号')}
                rules={[{ required: true, message: t('请输入商户号') }]}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='WechatApiV2Key'
                label={t('API v2 密钥')}
                placeholder={t('微信支付 API v2 密钥，敏感信息不显示')}
                type='password'
                rules={[{ required: true, message: t('请输入 API v2 密钥') }]}
              />
            </Col>
          </Row>
          <Row
            gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}
            style={{ marginTop: 16 }}
          >
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.InputNumber
                field='WechatMinTopUp'
                label={t('最低充值金额')}
                placeholder={t('例如：1，就是最低充值1元')}
              />
            </Col>
          </Row>
          <Button onClick={submitWechatSetting}>{t('更新微信支付设置')}</Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}
