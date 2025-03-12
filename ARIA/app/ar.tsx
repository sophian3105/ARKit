import { StyleSheet, Text, View } from "react-native";
import {
    ViroARScene,
    ViroARSceneNavigator,
    ViroText,
    ViroTrackingReason,
    ViroTrackingStateConstants,
} from "@reactvision/react-viro"
import { useState } from "react";

const TestSceneAR = () => {
    const [text, setText] = useState("Initializing Test AR..")

    function onInitialized(state: any, reason: ViroTrackingReason) {
        console.log("Test Initialized", state, reason);
        if (state == ViroTrackingStateConstants.TRACKING_NORMAL) {
            setText("Initialized and tracking!")
        }
    }

    return (
        <ViroARScene onTrackingUpdated={onInitialized}>
            <ViroText 
                text={text}
                scale={[0.5, 0.5, 0.5]}
                position={[0, 0, -1]}
                style={styles.testTextStyle}
            />
        </ViroARScene>
    )
}

export default function Index() {
  return (
    <ViroARSceneNavigator 
        autofocus={true}
        initialScene={{
            scene: TestSceneAR
        }}
        style={styles.index}
    />
  );
}

const styles = StyleSheet.create({
    index: { flex: 1 },
    testTextStyle: {
        fontSize: 30,
        color: "#ffffff",
        textAlignVertical: "center",
        textAlign: "center"
    }
});
