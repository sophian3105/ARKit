import { StyleSheet, Text, View, TouchableOpacity } from "react-native";
import { useRouter } from "expo-router"; // Use expo-router instead of react-navigation

export default function Index() {
  const router = useRouter(); // Get router object

  return (
    <View style={styles.container}>
      <Text style={styles.title}>ARIA</Text>
      <Text style={styles.description}>
        Record moments and follow lyrics in real-time. Relive the music anytime.
      </Text>
      
      {}
      <TouchableOpacity 
        style={styles.button} 
        onPress={() => router.push('/login')} 
      >
        <Text style={styles.buttonText}>Login</Text>
      </TouchableOpacity>

      {}
      <TouchableOpacity 
        style={styles.button} 
        onPress={() => router.push('/signup')} 
        >
            <Text style={styles.buttonText}>Sign Up</Text>
      </TouchableOpacity>

      {}
      <TouchableOpacity 
        style={styles.button}
        onPress={() => router.push('/ar')}
        >
            <Text style={styles.buttonText}>AR</Text>
        </TouchableOpacity>

      <Text style={styles.terms}>
        By continuing you agree to our Terms of Service and acknowledge our Privacy Policy.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'black',
    padding: 16,
  },
  title: {
    fontSize: 48,
    color: 'blue',
    marginBottom: 20,
  },
  description: {
    color: 'white',
    textAlign: 'center',
    marginBottom: 40,
  },
  button: {
    backgroundColor: 'blue',
    padding: 15,
    borderRadius: 5,
    width: '80%',
    alignItems: 'center',
    marginBottom: 20,
  },
  buttonText: {
    color: 'white',
    fontSize: 18,
  },
  terms: {
    color: 'white',
    textAlign: 'center',
    marginTop: 20,
  },
});
