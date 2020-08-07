import React from 'react';
import { Tabs, Spin, Row, Col, Radio } from 'antd';
import { Gallery } from './Gallery';
import { CreatePostButton } from "./CreatePostButton";
import { AroundMap } from './AroundMap';
import {
    GEOLOCATION_OPTIONS,
    POSITION_KEY,
    TOKEN_KEY,
    AUTH_HEADER,
    API_ROOT,
    POST_TYPE_IMAGE,
    POST_TYPE_VIDEO,
    TOPIC_AROUND,
    TOPIC_FACE,
} from "../constants";
import '../styles/Home.css';

const { TabPane } = Tabs;

export class Home extends React.Component {
    state = {
        loadingGeolocation: false,
        loadingPosts: false,
        errorMessage: null,
        posts: [],

        topic: TOPIC_AROUND,
    };

    onTopicChange = e => {
        console.log('radio checked', e.target.value);
        this.setState({
            topic: e.target.value,
        }, this.loadPost);
    };

    getGeolocation() {
        this.setState({
            loadingGeolocation: true,
            errorMessage: null,
        });

        //check if browser supports geolocaiton
        if ('geolocation' in navigator) {
            navigator.geolocation.getCurrentPosition(
                this.onGeolocationSuccess,
                this.onGeolocationFailure,
                GEOLOCATION_OPTIONS,
            );
        } else {
            this.setState({
                loadingGeolocation: false,
                errorMessage: 'Your browser does not support geolocation.',
            });
        }
    }

    onGeolocationSuccess = (position) => {
        this.setState({
            loadingGeolocation: false,
            errorMessage: null,
        });
        console.log(position);

        //destructuring
        const { latitude, longitude } = position.coords;
        //{ latitude: latitude, longitude: longitude }
        localStorage.setItem(POSITION_KEY, JSON.stringify({ latitude, longitude }));
        this.loadPost();
    }

    onGeolocationFailure = () => {
        this.setState({
            loadingGeolocation: false,
            errorMessage: 'Failed to load geolocation.',
        });
    }

    loadPost = (
        position = JSON.parse(localStorage.getItem(POSITION_KEY)),
        range = 20,
    ) => {
        if (this.state.topic === TOPIC_AROUND) {
            this.loadNearbyPost(position, range);
        } else if (this.state.topic === TOPIC_FACE) {
            this.loadFacePost();
        }
    }

    loadNearbyPost = (
        position,
        range
    ) => {
        this.setState({
            loadingPosts: true,
            errorMessage: null,
        });

        // const position = JSON.parse(localStorage.getItem(POSITION_KEY));
        // const range = 20;
        const token = localStorage.getItem(TOKEN_KEY);

        fetch(`${API_ROOT}/search?lat=${position.latitude}&lon=${position.longitude}&range=${range}`, {
            method: 'GET',
            headers: {
                Authorization: `${AUTH_HEADER} ${token}`,
            }
        }).then((response) => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Fail to load posts');
        }).then((data) => {
            console.log(data);
            this.setState({
                loadingPosts: false,
                posts: data ? data : [],
            });
        }).catch((error) => {
            this.setState({
                loadingPosts: false,
                errorMessage: error.message,
            });
        })
    }

    loadFacePost = () => {
        this.setState({
            loadingPosts: true,
            errorMessage: null,
        });

        const token = localStorage.getItem(TOKEN_KEY);

        fetch(`${API_ROOT}/cluster?term=face`, {
            method: 'GET',
            headers: {
                Authorization: `${AUTH_HEADER} ${token}`,
            }
        }).then((response) => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Fail to load posts');
        }).then((data) => {
            console.log(data);
            this.setState({
                loadingPosts: false,
                posts: data ? data : [],
            });
        }).catch((error) => {
            this.setState({
                loadingPosts: false,
                errorMessage: error.message,
            });
        })
    }

    getPosts(type) {
        if (this.state.errorMessage) {
            return <div>{this.state.errorMessage}</div>;
        } else if (this.state.loadingGeolocation) {
            return <Spin tip="Loading geolocation..."/>;
        } else if (this.state.loadingPosts) {
            return <Spin tip='Loading posts...'/>;
        } else if (this.state.posts.length > 0) {
            switch (type) {
                case POST_TYPE_IMAGE:
                    return this.getImagePosts();
                case POST_TYPE_VIDEO:
                    return this.getVideoPosts();
                default:
                    throw new Error('Unknown post type');
            }
        } else {
            return 'No nearby posts.';
        }
    }

    getImagePosts() {
        const images = this.state.posts
            .filter((post) => post.type === POST_TYPE_IMAGE)
            .map((post) => {
            return {
                src: post.url,
                thumbnail: post.url,
                thumbnailWidth: 400,
                thumbnailHeight: 300,
                caption: post.message,
                user: post.user,
            };
        });

        return (
            <Gallery images={images}/>
        );
    }

    getVideoPosts() {
        return (
            <Row gutter={30}>
                {
                    this.state.posts
                        .filter((post) => post.type === POST_TYPE_VIDEO)
                        .map((post) => (
                            <Col className="gutter-row" span={6}>
                                <video src={post.url} controls className="video-block" />
                                <div>{`${post.user}: ${post.message}`}</div>
                            </Col>
                        ))
                }
            </Row>
        );
    }

    componentDidMount() {
        this.getGeolocation();
    }

    render() {
        const operations = <CreatePostButton onSuccess={this.loadPost} />;

        return (
            <div>
                <Radio.Group onChange={this.onTopicChange} value={this.state.topic} className="topic-radio-group">
                    <Radio value={TOPIC_AROUND}>Posts Around Me</Radio>
                    <Radio value={TOPIC_FACE}>Faces Around the World</Radio>
                </Radio.Group>

                <Tabs tabBarExtraContent={operations} className="main-tabs">
                    <TabPane tab="Image Posts" key="1">
                        {this.getPosts(POST_TYPE_IMAGE)}
                    </TabPane>
                    <TabPane tab="Video Posts" key="2">
                        {this.getPosts(POST_TYPE_VIDEO)}
                    </TabPane>
                    <TabPane tab="Map" key="3">
                        <AroundMap
                            googleMapURL="https://maps.googleapis.com/maps/api/js?key=AIzaSyD3CEh9DXuyjozqptVB5LA-dN7MxWWkr9s&v=3.exp&libraries=geometry,drawing,places"
                            loadingElement={<div style={{ height: `100%` }} />}
                            containerElement={<div style={{ height: `600px` }} />}
                            mapElement={<div style={{ height: `100%` }} />}
                            posts={this.state.posts}
                            onChange={this.loadPost}
                        />
                    </TabPane>
                </Tabs>
            </div>

        );
    }
};