import React from 'react';
import { Form, Input, Upload } from 'antd';
import { InboxOutlined } from '@ant-design/icons';

class NormalCreatePostForm extends React.Component {
    normFile = e => {
        console.log('Upload event:', e);
        if (Array.isArray(e)) {
            return e;
        }
        return e && e.fileList;
    };

    render() {
        // const { getFieldDecorator } = this.props.form;

        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 14 },
        };

        return (
            <Form
                // name="validate_other"
                {...formItemLayout}
                // onFinish={onFinish}
                // initialValues={{
                //     ['input-number']: 3,
                //     ['checkbox-group']: ['A', 'B'],
                //     rate: 3.5,
                // }}
            >
                <Form.Item label="Message" name="message"
                           rules={[
                               {
                                   required: true,
                                   message: 'Please input your message!',
                               },
                           ]}>
                    <Input />
                </Form.Item>

                <Form.Item label="Image">
                    <Form.Item name="image" valuePropName="fileList" getValueFromEvent={this.normFile} noStyle
                       rules={[
                        {
                            required: true,
                            message: 'Please select an image!',
                        },
                    ]}>
                        <Upload.Dragger name="files" beforeUpload={() => false}>
                            <p className="ant-upload-drag-icon">
                                <InboxOutlined/>
                            </p>
                            <p className="ant-upload-text">Click or drag file to this area to upload</p>
                            <p className="ant-upload-hint">Support for a single or bulk upload.</p>
                        </Upload.Dragger>
                    </Form.Item>
                </Form.Item>

            </Form>
        );
    }
}

export const CreatePostForm = NormalCreatePostForm;

