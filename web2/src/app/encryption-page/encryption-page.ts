import { Component, computed, inject, signal } from '@angular/core';
import { EncryptionForm } from '../encryption-form';
import { form, FormField, minLength, required, validate } from '@angular/forms/signals';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-encryption-page',
  imports: [FormField],
  templateUrl: './encryption-page.html',
  styleUrl: './encryption-page.css',
})
export class EncryptionPage {
  private http = inject(HttpClient);

  private formData = signal<EncryptionForm>({
    password1: "",
    password2: "",
    name: "",
    data: "",
    output: "",
  });

  encForm = form(this.formData, (schemaPath) => {
    required(schemaPath.password1, {message: "password is required"})
    minLength(schemaPath.password1, 8, {message: "password must be at least 8 characters"})
    required(schemaPath.name, {message: "Must provide credential name"})
    required(schemaPath.data, {message: "Must provide data to encrypt/decrypt"})

    validate(schemaPath.password2, ({value, valueOf}) => {
      const confirmPassword = value();
      const password = valueOf(schemaPath.password1);
      if (confirmPassword !== password) {
        return {
          kind: 'passwordMismatch',
          message: 'Passwords do not match',
        };
      }
      return null;
    });
  });

  public encrypt() {
    console.log('encrypt');

    let postData = {
      password: this.encForm.password1().value(),
      data: this.encForm.data().value(),
      name: this.encForm.name().value(),
    }

    this.http.post('/encrypt', postData).subscribe((result: any) => {
      if (result.success) {
        this.encForm.output().value.set(result.result);
      } else {
        this.encForm.output().value.set('Error: ' + result.message);
      }
    }, (err) => {
      this.encForm.output().value.set(err.message);
    });
  }

  public decrypt() {
    console.log('decrypt');

    let postData = {
      password: this.encForm.password1().value(),
      data: this.encForm.data().value(),
      name: this.encForm.name().value(),
    }

    this.http.post('/decrypt', postData).subscribe((result: any) => {
      if (result.success) {
        this.encForm.output().value.set(result.result);
      } else {
        this.encForm.output().value.set('Error: ' + result.message);
      }
    }, (err) => {
      this.encForm.output().value.set(err.message);
    });
  }

  public copy() {
    navigator.clipboard.writeText(this.encForm.output().value() ?? '');
  }
}
