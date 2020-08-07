import React from 'react';
import { Link } from 'react-router-dom';
import {API_ROOT} from "../constants"
import '../styles/Login.css';

import { Form, Input, Button, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

class NormalLoginForm extends React.Component {
    handleSubmit = e => {
        e.preventDefault();
        let lastResponse;

        this.props.form.validateFields((err, values) => {
            if (!err) {
                console.log('Received values of form: ', values);
                fetch(`${API_ROOT}/login`, {
                    method: 'POST',
                    body: JSON.stringify({
                        username: values.username,
                        password: values.password,
                    }),
                }).then((response) => {
                    lastResponse = response;
                    return response.text();
                }, (error) => {
                    console.log('Error');
                }).then((text) => {
                    if (lastResponse.ok) {
                        message.success('Login success!');
                        //successful login => return token (text)
                        this.props.handleLogin(text);
                    } else {
                        message.error(text);
                    }
                });
            }
        });
    };

    render () {
        return (
            <Form onSubmit={this.handleSubmit} className="login-form">
                <Form.Item name="username"
                   rules={[
                       {
                           required: true,
                           message: 'Please input your username!',
                       },
                   ]}>
                    <Input
                        prefix={<UserOutlined className="site-form-item-icon"/>} placeholder="Username"
                        // placeholder="Username"
                    />
                </Form.Item>

                <Form.Item name="password"
                       rules={[
                            {
                            required: true,
                            message: 'Please input your password!',
                            },
                ]}>
                        <Input
                            prefix={<LockOutlined className="site-form-item-icon"/>} type="password" placeholder="Password"
                            // type="password"
                            // placeholder="Password"
                        />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" className="login-form-button">
                        Log in
                    </Button>
                    <div>
                        Or <Link to="/register">register now!</Link>
                    </div>
                </Form.Item>
            </Form>
        );
    }

    // render () {
    //     const { getFieldDecorator } = this.props.form;
    //
    //     return (
    //         <Form onSubmit={this.handleSubmit} className="login-form">
    //             <Form.Item>
    //                 {getFieldDecorator('username', {
    //                     rules: [
    //                         {
    //                             required: true,
    //                             message: 'Please input your username!',
    //                         },
    //                     ],
    //                 })(
    //                     <Input
    //                       prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />}
    //                           placeholder="Username"
    //                       />,
    //                   )}
    //             </Form.Item>
    //
    //             <Form.Item>
    //                 {getFieldDecorator('password', {
    //                     rules: [
    //                         {
    //                             required: true,
    //                             message: 'Please input your password!',
    //                         },
    //                         {
    //                             validator: this.validateToNextPassword,
    //                         }
    //                     ],
    //                 })(
    //                     <Input
    //                         prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />}
    //                         type="password"
    //                         placeholder="Password"
    //                     />,
    //                 )}
    //             </Form.Item>
    //
    //             <Form.Item>
    //                 <Button type="primary" htmlType="submit" className="login-form-button">
    //                     Log in
    //                 </Button>
    //                 <div>
    //                     Or <Link to="/register">register now!</Link>
    //                 </div>
    //             </Form.Item>
    //         </Form>
    //     );
    // }
};

export const Login = NormalLoginForm;
// export const Login = Form.create({ name: 'normal_login' })(NormalLoginForm);