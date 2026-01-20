/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React, { useEffect, useRef } from 'react';
import { Modal, Typography, Spin } from '@douyinfe/semi-ui';
import { QRCodeSVG } from 'qrcode.react';
import { SiAlipay, SiWechat } from 'react-icons/si';
import { API, showSuccess } from '../../../helpers';

const { Text, Title } = Typography;

const QRCodePaymentModal = ({
  t,
  visible,
  onCancel,
  payUrl,
  payWay,
  amount,
  orderId,
  onPaymentSuccess,
}) => {
  const pollingTimerRef = useRef(null);
  const isWechat = payWay === 'wechat' || payWay === 'wxpay';
  const isAlipay = payWay === 'alipay';

  // 轮询订单状态
  const checkOrderStatus = async () => {
    if (!orderId) return;

    try {
      const res = await API.get(`/api/user/topup/status?order_id=${orderId}`);
      const { message, data } = res.data;

      if (message === 'success' && data.status === 'success') {
        // 支付成功
        showSuccess(t('支付成功！'));
        if (onPaymentSuccess) {
          onPaymentSuccess();
        }
        stopPolling();
        onCancel();
      }
    } catch (error) {
      console.error('查询订单状态失败:', error);
    }
  };

  const startPolling = () => {
    if (pollingTimerRef.current) {
      clearInterval(pollingTimerRef.current);
    }
    // 每2秒检查一次订单状态
    pollingTimerRef.current = setInterval(checkOrderStatus, 2000);
  };

  const stopPolling = () => {
    if (pollingTimerRef.current) {
      clearInterval(pollingTimerRef.current);
      pollingTimerRef.current = null;
    }
  };

  useEffect(() => {
    if (visible && orderId) {
      startPolling();
    } else {
      stopPolling();
    }

    return () => stopPolling();
  }, [visible, orderId]);

  const handleCancel = () => {
    stopPolling();
    onCancel();
  };

  const getTitle = () => {
    if (isWechat) return t('微信支付');
    if (isAlipay) return t('支付宝支付');
    return t('扫码支付');
  };

  const getIcon = () => {
    if (isWechat) return <SiWechat size={24} color="#07C160" />;
    if (isAlipay) return <SiAlipay size={24} color="#1677FF" />;
    return null;
  };

  const getBgColor = () => {
    if (isWechat) return '#07C160';
    if (isAlipay) return '#1677FF';
    return '#666';
  };

  return (
    <Modal
      title={
        <div className="flex items-center gap-2">
          {getIcon()}
          <span>{getTitle()}</span>
        </div>
      }
      visible={visible}
      onCancel={handleCancel}
      footer={null}
      maskClosable={false}
      centered
      width={400}
    >
      <div className="flex flex-col items-center py-4">
        {payUrl ? (
          <>
            <div
              className="p-4 rounded-xl mb-4"
              style={{ backgroundColor: getBgColor() + '10' }}
            >
              <QRCodeSVG
                value={payUrl}
                size={200}
                level="H"
                includeMargin={true}
                bgColor="#ffffff"
                fgColor="#000000"
              />
            </div>
            <Text strong className="text-lg mb-2">
              {t('请使用')}
              {isWechat ? t('微信') : isAlipay ? t('支付宝') : ''}
              {t('扫码支付')}
            </Text>
            {amount && (
              <Text className="text-2xl font-bold mb-2" style={{ color: getBgColor() }}>
                ¥{amount}
              </Text>
            )}
            <Text type="tertiary" size="small">
              {t('支付完成后会自动刷新')}
            </Text>
            {orderId && (
              <Text type="tertiary" size="small" className="mt-2">
                {t('订单号')}: {orderId}
              </Text>
            )}
          </>
        ) : (
          <Spin size="large" tip={t('正在生成支付二维码...')} />
        )}
      </div>
    </Modal>
  );
};

export default QRCodePaymentModal;
