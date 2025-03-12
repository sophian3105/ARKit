import React from 'react';
import { View, TextInput, Button, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useRouter } from 'expo-router';

export default function Login() {
  const router = useRouter();

  return (
    <View style={styles.container}>
      <TextInput style={styles.input} placeholder="Email" placeholderTextColor="black" />
      <TextInput style={styles.input} placeholder="Password" placeholderTextColor="black" secureTextEntry />
      
      <TouchableOpacity style={styles.button} onPress={() => {}}>
        <Text style={styles.buttonText}>Login</Text>
      </TouchableOpacity>

      <Text style={styles.or}>Or login with</Text>

      {/* Third-party apps */}

      <TouchableOpacity onPress={() => router.push('/signup')}>
        <Text style={styles.loginLink}>
          Don't have an account? <Text style={styles.signupText}>Sign up here</Text>
        </Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 16,
    backgroundColor: 'black',
  },
  input: {
    height: 50,
    backgroundColor: 'white',
    borderWidth: 1,
    marginBottom: 12,
    paddingHorizontal: 10,
    color: 'black',
    borderRadius: 8,
  },
  button: {
    backgroundColor: 'blue', 
    padding: 15,
    alignItems: 'center',
    borderRadius: 8,
    marginVertical: 10,
  },
  buttonText: {
    color: 'white',
    fontSize: 18,
  },
  or: {
    textAlign: 'center',
    color: 'white',
    marginVertical: 10,
  },
  loginLink: {
    textAlign: 'center',
    color: 'blue',
    marginTop: 10,
  },
  signupText: {
    color: 'white',
  },
});
