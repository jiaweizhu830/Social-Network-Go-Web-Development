import React from 'react';
import { Modal, Button, message } from 'antd';
import { CreatePostForm } from "./CreatePostForm";
import { API_ROOT, AUTH_HEADER, POSITION_KEY, POSITION_NOISE } from "../constants"

export class CreatePostButton extends React.Component {
    state = {
        visible: false,
        confirmLoading: false,
    };

    showModal = () => {
        this.setState({
            visible: true,
        });
    };

    handleOk = () => {
        this.setState({
            confirmLoading: true,
        });

        this.form.scrollToField((err, values) => {
            if (!err) {
                console.log('Received values of form: ', values);

                //上传FormData, 因为不允许直接上传文件
                const formData = new FormData();
                const token = localStorage.getItem(POSITION_KEY);
                const position = JSON.parse(localStorage.getItem(POSITION_KEY));

                //add noise to position
                formData.append('lat', position.latitude + Math.random() * POSITION_NOISE * 2 - POSITION_NOISE);
                formData.append('lon', position.longitude + Math.random() * POSITION_NOISE * 2 - POSITION_NOISE);
                formData.append('message', values.message);
                //只上传第一张图片 if 多图片上传
                formData.append('image', values.image[0].originFileObj);

                fetch(`${API_ROOT}/post`, {
                    method: 'POST',
                    body: formData,
                    headers: {
                        Authorization: `${AUTH_HEADER} ${token}`,
                    },
                    dataType: 'text',
                }).then((response) => {
                    if (response.ok) {
                        message.success('Create post succeeded!');
                        this.form.resetFields();
                        this.setState({
                            visible: false,
                            confirmLoading: false,
                        });

                        if (this.props.onSuccess) {
                            this.props.onSuccess();
                        }
                    } else {
                        message.error('Create  post failed.');
                        this.setState({
                            confirmLoading: false,
                        });
                    }
                })
            } else {
                this.setState({
                    confirmLoading: false,
                });
            }
        })
    };

    handleCancel = () => {
        console.log('Clicked cancel button');
        this.setState({
            visible: false,
        });
    };

    saveFormRef = (formInstance) => {
        this.form = formInstance;
    }

    render() {
        const { visible, confirmLoading } = this.state;
        return (
            <div>
                <Button type="primary" onClick={this.showModal}>
                    Create New Post
                </Button>
                <Modal
                    title="Create New Post"
                    okText="Create"
                    visible={visible}
                    onOk={this.handleOk}
                    confirmLoading={confirmLoading}
                    onCancel={this.handleCancel}
                >
                    <CreatePostForm ref={this.saveFormRef} />
                </Modal>
            </div>
        );
    }
}
