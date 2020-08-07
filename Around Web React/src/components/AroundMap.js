import React from "react";
import {
    withScriptjs,
    withGoogleMap,
    GoogleMap,
} from "react-google-maps"
import { AroundMarker } from "./AroundMarker"
import { POSITION_KEY } from "../constants"

export class NormalAroundMap extends React.Component {
    reloadMarkers = () => {
        const center = this.map.getCenter();
        const position = {
            latitude: center.lat(),
            longitude: center.lng(),
        };

        //map 矩形
        const bounds = this.map.getBounds();
        const northEast = bounds.getNorthEast();
        //正东边界点: 右上角竖线 与 中心横线的 交点
        const east = new window.google.maps.LatLng(center.lat(), northEast.lng());
        //compute 球体上两点距离 (center 与map 矩形正东的边界点)
        //convert m => km
        const range =
            window.google.maps.geometry.spherical.computeDistanceBetween(center, east)
            / 1000;

        //loadNearby posts
        //Check is this prop is assigned
        if (this.props.onChange) {
            this.props.onChange(position, range);
        }
    }

    saveMapRef = (mapInstance) => {
        this.map = mapInstance;
        window.map = mapInstance;
    }

    render() {
        const position = JSON.parse(localStorage.getItem(POSITION_KEY));

        return(
            <GoogleMap
                ref={this.saveMapRef}
                defaultZoom={11}
                defaultCenter={{ lat: position.latitude, lng: position.longitude }}
                onDragEnd={this.reloadMarkers}
                onZoomChanged={this.reloadMarkers}
                onResize={this.reloadMarkers}
            >
                {
                    this.props.posts.map((post) => (
                        <AroundMarker
                            post={post}
                            key={post.url}
                        />
                    ))
                }
            </GoogleMap>
        );
    }
}

export const AroundMap = withScriptjs(withGoogleMap(NormalAroundMap));