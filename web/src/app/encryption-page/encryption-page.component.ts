import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import {
  AbstractControl,
  FormControl,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';

@Component({
  selector: 'app-encryption-page',
  templateUrl: './encryption-page.component.html',
  styleUrls: ['./encryption-page.component.scss'],
})
export class EncryptionPageComponent {
  constructor(private http: HttpClient) {}

  form = new FormGroup(
    {
      password1: new FormControl('', [
        Validators.required,
        Validators.minLength(8),
      ]),
      password2: new FormControl('', [
        Validators.required,
        Validators.minLength(8),
      ]),
      name: new FormControl('', [Validators.required]),
      data: new FormControl('', [Validators.required]),
      output: new FormControl(''),
    },
    { validators: passwordsMatchValidator }
  );

  public failedSubmit = false;

  public encrypt() {
    console.log('encrypt');
    if (this.form.invalid) {
      this.failedSubmit = true;
      console.log('invalid');
      return;
    }

    this.failedSubmit = false;

    let postData = {
      password: this.form.value.password1,
      data: this.form.value.data,
      name: this.form.value.name
    }

    this.http.post('/encrypt', postData).subscribe((result: any) => {
      if (result.success) {
        this.form.patchValue({ output: result.result });
      } else {
        this.form.patchValue({ output: 'Error: ' + result.message });
      }
    });
  }

  public decrypt() {
    console.log('decrypt');
    if (this.form.invalid) {
      this.failedSubmit = true;
      console.log('invalid');
      return;
    }

    this.failedSubmit = false;

    let postData = {
      password: this.form.value.password1,
      data: this.form.value.data,
      name: this.form.value.name
    }

    this.http.post('/decrypt', postData).subscribe((result: any) => {
      if (result.success) {
        this.form.patchValue({ output: result.result });
      } else {
        this.form.patchValue({ output: 'Error: ' + result.message });
      }
    });
  }
}

const passwordsMatchValidator: ValidatorFn = (
  control: AbstractControl
): ValidationErrors | null => {
  const pass1 = control.get('password1');
  const pass2 = control.get('password2');

  return pass1 && pass2 && pass1.value === pass2.value
    ? null
    : { passwordsDoNotMatch: true };
};
